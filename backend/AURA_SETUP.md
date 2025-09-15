# Neo4j Aura Setup Guide

This guide will help you set up your NQ application with Neo4j Aura cloud database.

## 1. Create Neo4j Aura Instance

1. Go to [Neo4j Aura Console](https://console.neo4j.io/)
2. Sign up or log in to your account
3. Click "New Instance"
4. Choose your plan (Free tier is perfect for development)
5. Select a region close to your users
6. Give your database a name (e.g., "nq-media-tracker")

## 2. Get Connection Details

After creating your instance:

1. Click on your database instance
2. Go to the "Connect" tab
3. Copy the following details:
   - **Connection URI**: Looks like `neo4j+s://abcd1234.databases.neo4j.io`
   - **Username**: Usually `neo4j`
   - **Password**: The password you set when creating the instance

## 3. Configure Environment Variables

1. Copy `env.example` to `.env`:
   ```bash
   cp env.example .env
   ```

2. Update your `.env` file with your Aura credentials:
   ```bash
   # Neo4j Aura Configuration
   NEO4J_URI=neo4j+s://your-instance-id.databases.neo4j.io
   NEO4J_USERNAME=neo4j
   NEO4J_PASSWORD=your-aura-password
   
   # Server Configuration
   PORT=8080
   ```

## 4. Test the Connection

Run your application to test the connection:

```bash
go run .
```

You should see:
```
Successfully connected to Neo4j Aura
connect to http://localhost:8080/ for GraphQL playground
```

## 5. Access Neo4j Browser (Optional)

You can also access your database through the Neo4j Browser:

1. In your Aura console, click "Open with Neo4j Browser"
2. Use the same credentials to explore your data
3. Run Cypher queries to see your data

## Troubleshooting

### Connection Issues
- **Wrong URI format**: Make sure your URI starts with `neo4j+s://` and contains `databases.neo4j.io`
- **Authentication failed**: Double-check your username and password
- **Network issues**: Ensure your firewall allows outbound HTTPS connections

### Performance Tips
- The free tier has limited resources, so avoid running heavy queries
- Use the connection pooling settings optimized for Aura (already configured)
- Monitor your usage in the Aura console

### Common Aura URI Formats
```
# Secure connection (recommended)
neo4j+s://instance-id.databases.neo4j.io

# Non-secure connection (not recommended for production)
neo4j://instance-id.databases.neo4j.io
```

## Next Steps

1. **Start using your GraphQL API**: Visit `http://localhost:8080` to access the playground
2. **Create some test data**: Try creating users and media items
3. **Explore your data**: Use the Neo4j Browser to see your data visually
4. **Monitor usage**: Keep an eye on your Aura console for usage metrics

## Security Notes

- Never commit your `.env` file to version control
- Use strong passwords for your Aura instances
- Consider using Aura's IP whitelisting for production environments
- The `neo4j+s://` protocol provides encrypted connections
