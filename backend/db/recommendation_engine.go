package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// BuildRecommendations computes and upserts backend recommendations for the given user.
//
// Scoring uses two complementary signals:
//  1. Tag/graph similarity  — media connected via RELATED_TO to things the user has
//     already logged receives a score proportional to how many direct RELATED_TO hops
//     exist (closer == higher score, up to depth 2).
//  2. Friends boost         — media that one or more friends have logged (but the user
//     has not) receives +10 score per friend who logged it.
//
// Results are upserted as Recommendation nodes with source='engine', deduplicated by
// (userId, mediaId). Existing engine recommendations for other media are cleared first
// so stale items don't linger.
func (r *Neo4jRepository) BuildRecommendations(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Step 1: Delete old engine recommendations for this user so we start fresh.
		_, err := tx.Run(ctx, `
			MATCH (u:User {id: $userID})-[:RECEIVED_RECOMMENDATION]->(rec:Recommendation {source: 'engine'})
			DETACH DELETE rec
		`, map[string]any{"userID": userID.String()})
		if err != nil {
			return nil, err
		}

		// Step 2: Tag/graph similarity score.
		// For each media the user has logged, traverse RELATED_TO up to 2 hops.
		// Score = 20 for depth-1 neighbours, 10 for depth-2 neighbours (before aggregation).
		// We aggregate by candidate media and sum scores.
		_, err = tx.Run(ctx, `
			MATCH (u:User {id: $userID})-[:HAS_ACTIVITY]->(seed:Media)
			MATCH p = (seed)-[:RELATED_TO*1..2]-(candidate:Media)
			WHERE NOT (u)-[:HAS_ACTIVITY]->(candidate)
			  AND candidate <> seed
			WITH u, candidate,
			     sum(CASE length(p) WHEN 1 THEN 20 WHEN 2 THEN 10 ELSE 5 END) AS graphScore
			WHERE graphScore > 0
			MERGE (rec:Recommendation {userId: $userID, mediaId: candidate.id, source: 'engine'})
			ON CREATE SET rec.id        = randomUUID(),
			              rec.score     = graphScore,
			              rec.createdAt = datetime()
			ON MATCH  SET rec.score     = graphScore
			WITH u, rec, candidate
			MERGE (u)-[:RECEIVED_RECOMMENDATION]->(rec)
			MERGE (rec)-[:RECOMMENDS]->(candidate)
		`, map[string]any{"userID": userID.String()})
		if err != nil {
			return nil, err
		}

		// Step 3: Friends boost.
		// For each friend's logged media not yet logged by the user, add +10 per friend.
		// Uses MERGE so if the rec already exists from step 2 we increment its score.
		_, err = tx.Run(ctx, `
			MATCH (u:User {id: $userID})-[:FRIEND_OF]->(friend:User)
			MATCH (friend)-[:HAS_ACTIVITY]->(candidate:Media)
			WHERE NOT (u)-[:HAS_ACTIVITY]->(candidate)
			WITH u, candidate, count(DISTINCT friend) AS friendCount
			WHERE friendCount > 0
			MERGE (rec:Recommendation {userId: $userID, mediaId: candidate.id, source: 'engine'})
			ON CREATE SET rec.id        = randomUUID(),
			              rec.score     = toFloat(friendCount) * 10.0,
			              rec.createdAt = datetime()
			ON MATCH  SET rec.score     = rec.score + toFloat(friendCount) * 10.0
			WITH u, rec, candidate
			MERGE (u)-[:RECEIVED_RECOMMENDATION]->(rec)
			MERGE (rec)-[:RECOMMENDS]->(candidate)
		`, map[string]any{"userID": userID.String()})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	return err
}
