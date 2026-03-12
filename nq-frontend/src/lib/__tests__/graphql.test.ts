import { 
  ME_QUERY, 
  ME_ACTIVITIES_QUERY,
  GET_HOME_MEDIA_QUERY,
  GET_MEDIA_DETAILS_QUERY,
  AUTOCOMPLETE_MEDIA_QUERY,
  RECURSIVE_SEARCH_STATUS_QUERY,
  UPDATE_USER_MUTATION,
} from '../graphql';

describe('GraphQL Queries and Mutations', () => {
  describe('Query exports', () => {
    it('should export ME_QUERY as a DocumentNode', () => {
      expect(ME_QUERY).toBeDefined();
      expect(ME_QUERY.kind).toBe('Document');
    });

    it('should export ME_ACTIVITIES_QUERY as a DocumentNode', () => {
      expect(ME_ACTIVITIES_QUERY).toBeDefined();
      expect(ME_ACTIVITIES_QUERY.kind).toBe('Document');
    });

    it('should export GET_HOME_MEDIA_QUERY as a DocumentNode', () => {
      expect(GET_HOME_MEDIA_QUERY).toBeDefined();
      expect(GET_HOME_MEDIA_QUERY.kind).toBe('Document');
    });

    it('should export GET_MEDIA_DETAILS_QUERY as a DocumentNode', () => {
      expect(GET_MEDIA_DETAILS_QUERY).toBeDefined();
      expect(GET_MEDIA_DETAILS_QUERY.kind).toBe('Document');
    });

    it('should export AUTOCOMPLETE_MEDIA_QUERY as a DocumentNode', () => {
      expect(AUTOCOMPLETE_MEDIA_QUERY).toBeDefined();
      expect(AUTOCOMPLETE_MEDIA_QUERY.kind).toBe('Document');
    });

    it('should export RECURSIVE_SEARCH_STATUS_QUERY as a DocumentNode', () => {
      expect(RECURSIVE_SEARCH_STATUS_QUERY).toBeDefined();
      expect(RECURSIVE_SEARCH_STATUS_QUERY.kind).toBe('Document');
    });
  });

  describe('Mutation exports', () => {
    it('should export UPDATE_USER_MUTATION as a DocumentNode', () => {
      expect(UPDATE_USER_MUTATION).toBeDefined();
      expect(UPDATE_USER_MUTATION.kind).toBe('Document');
    });
  });

  describe('Query structure', () => {
    it('ME_QUERY should contain expected fields', () => {
      const queryString = ME_QUERY.loc?.source.body || '';
      expect(queryString).toContain('query Me');
      expect(queryString).toContain('me');
      expect(queryString).toContain('id');
      expect(queryString).toContain('name');
      expect(queryString).toContain('email');
    });
  });
});
