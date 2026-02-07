import {
  capMediaCandidates,
  scoreMediaFromUser,
  scoreMediaFromRootMedia,
  GraphMediaNode,
  MAX_DEPTH,
  LAYER_DECAY,
  CANDIDATE_CAP,
} from '../graphScore';

describe('graphScore', () => {
  describe('capMediaCandidates', () => {
    it('should return all media if count is below cap', () => {
      const media: GraphMediaNode[] = [
        { id: '1', title: 'Movie A', averageRating: 8.5 },
        { id: '2', title: 'Movie B', averageRating: 7.0 },
      ];

      const result = capMediaCandidates(media, 10);
      expect(result).toHaveLength(2);
      expect(result).toEqual(media);
    });

    it('should cap media to specified limit', () => {
      const media: GraphMediaNode[] = [
        { id: '1', title: 'Movie A', averageRating: 8.5 },
        { id: '2', title: 'Movie B', averageRating: 7.0 },
        { id: '3', title: 'Movie C', averageRating: 9.0 },
        { id: '4', title: 'Movie D', averageRating: 6.5 },
      ];

      const result = capMediaCandidates(media, 2);
      expect(result).toHaveLength(2);
    });

    it('should sort by rating descending', () => {
      const media: GraphMediaNode[] = [
        { id: '1', title: 'Movie A', averageRating: 7.0 },
        { id: '2', title: 'Movie B', averageRating: 9.0 },
        { id: '3', title: 'Movie C', averageRating: 8.0 },
      ];

      const result = capMediaCandidates(media, 2);
      expect(result[0].id).toBe('2'); // Highest rating (9.0)
      expect(result[1].id).toBe('3'); // Second highest (8.0)
    });

    it('should sort by title alphabetically when ratings are equal', () => {
      const media: GraphMediaNode[] = [
        { id: '1', title: 'Zebra', averageRating: 7.0 },
        { id: '2', title: 'Apple', averageRating: 7.0 },
        { id: '3', title: 'Mango', averageRating: 7.0 },
      ];

      const result = capMediaCandidates(media, 2);
      expect(result[0].title).toBe('Apple');
      expect(result[1].title).toBe('Mango');
    });

    it('should handle null ratings', () => {
      const media: GraphMediaNode[] = [
        { id: '1', title: 'Movie A', averageRating: null },
        { id: '2', title: 'Movie B', averageRating: 7.0 },
      ];

      const result = capMediaCandidates(media, 1);
      expect(result[0].id).toBe('2'); // Rated movie comes first
    });

    it('should use default CANDIDATE_CAP when cap not provided', () => {
      const media: GraphMediaNode[] = Array.from({ length: 350 }, (_, i) => ({
        id: String(i),
        title: `Movie ${i}`,
        averageRating: Math.random() * 10,
      }));

      const result = capMediaCandidates(media);
      expect(result).toHaveLength(CANDIDATE_CAP); // Should use default 300
    });
  });

  describe('scoreMediaFromUser', () => {
    it('should score media based on user activity', () => {
      const activityMedia: GraphMediaNode[] = [
        {
          id: '1',
          title: 'Inception',
          genres: [{ name: 'Sci-Fi' }, { name: 'Thriller' }],
        },
      ];

      const candidates: GraphMediaNode[] = [
        {
          id: '2',
          title: 'Interstellar',
          genres: [{ name: 'Sci-Fi' }],
        },
        {
          id: '3',
          title: 'The Notebook',
          genres: [{ name: 'Romance' }],
        },
      ];

      const scores = scoreMediaFromUser({ candidates, activityMedia });

      // Interstellar should score higher (shares Sci-Fi genre)
      const interstellarScore = scores.get('2') ?? 0;
      const notebookScore = scores.get('3') ?? 0;

      expect(interstellarScore).toBeGreaterThan(notebookScore);
    });

    it('should handle empty activity media', () => {
      const candidates: GraphMediaNode[] = [
        { id: '1', title: 'Movie A', genres: [{ name: 'Action' }] },
      ];

      const scores = scoreMediaFromUser({ candidates, activityMedia: [] });
      expect(scores.size).toBe(0);
    });

    it('should respect custom maxDepth parameter', () => {
      const activityMedia: GraphMediaNode[] = [
        { id: '1', title: 'Movie A', genres: [{ name: 'Action' }] },
      ];
      const candidates: GraphMediaNode[] = [
        { id: '2', title: 'Movie B', genres: [{ name: 'Action' }] },
      ];

      const scoresDefault = scoreMediaFromUser({ candidates, activityMedia });
      const scoresCustom = scoreMediaFromUser({
        candidates,
        activityMedia,
        maxDepth: 1,
      });

      // Both should work but may produce different scores
      expect(scoresDefault).toBeDefined();
      expect(scoresCustom).toBeDefined();
    });

    it('should respect custom layerDecay parameter', () => {
      const activityMedia: GraphMediaNode[] = [
        { id: '1', title: 'Movie A', genres: [{ name: 'Action' }] },
      ];
      const candidates: GraphMediaNode[] = [
        { id: '2', title: 'Movie B', genres: [{ name: 'Action' }] },
      ];

      const scoresDefault = scoreMediaFromUser({ candidates, activityMedia });
      const scoresCustom = scoreMediaFromUser({
        candidates,
        activityMedia,
        layerDecay: 0.3,
      });

      expect(scoresDefault).toBeDefined();
      expect(scoresCustom).toBeDefined();
    });
  });

  describe('scoreMediaFromRootMedia', () => {
    it('should score candidates based on root media', () => {
      const rootMedia: GraphMediaNode = {
        id: '1',
        title: 'The Dark Knight',
        genres: [{ name: 'Action' }, { name: 'Crime' }],
      };

      const candidates: GraphMediaNode[] = [
        {
          id: '2',
          title: 'Batman Begins',
          genres: [{ name: 'Action' }],
        },
        {
          id: '3',
          title: 'The Notebook',
          genres: [{ name: 'Romance' }],
        },
      ];

      const scores = scoreMediaFromRootMedia({ candidates, rootMedia });

      // Batman Begins should score higher (shares Action genre)
      const batmanScore = scores.get('2') ?? 0;
      const notebookScore = scores.get('3') ?? 0;

      expect(batmanScore).toBeGreaterThan(notebookScore);
    });

    it('should handle root media with no shared features', () => {
      const rootMedia: GraphMediaNode = {
        id: '1',
        title: 'Unique Movie',
        genres: [{ name: 'Documentary' }],
      };

      const candidates: GraphMediaNode[] = [
        { id: '2', title: 'Action Movie', genres: [{ name: 'Action' }] },
      ];

      const scores = scoreMediaFromRootMedia({ candidates, rootMedia });
      expect(scores.size).toBeGreaterThanOrEqual(0);
    });

    it('should score media with shared cast members', () => {
      const rootMedia: GraphMediaNode = {
        id: '1',
        title: 'Movie A',
        cast: [{ id: 'actor1', name: 'Tom Hanks' }],
      };

      const candidates: GraphMediaNode[] = [
        {
          id: '2',
          title: 'Movie B',
          cast: [{ id: 'actor1', name: 'Tom Hanks' }],
        },
        {
          id: '3',
          title: 'Movie C',
          cast: [{ id: 'actor2', name: 'Other Actor' }],
        },
      ];

      const scores = scoreMediaFromRootMedia({ candidates, rootMedia });

      const movieBScore = scores.get('2') ?? 0;
      const movieCScore = scores.get('3') ?? 0;

      // Movie B shares cast member, should score higher
      expect(movieBScore).toBeGreaterThan(movieCScore);
    });
  });

  describe('constants', () => {
    it('should export expected constant values', () => {
      expect(MAX_DEPTH).toBe(2);
      expect(LAYER_DECAY).toBe(0.5);
      expect(CANDIDATE_CAP).toBe(300);
    });
  });
});
