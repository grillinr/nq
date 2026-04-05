# NQ Backend Production Deployment

This directory contains everything needed to deploy the NQ GraphQL backend as a production-ready containerized service accessible at `https://nq.nathangrilliot.com`.

## 🚀 Quick Start

1. **Set up firewall** (requires sudo):
   ```bash
   sudo ./scripts/setup-firewall.sh
   ```

2. **Deploy the production stack**:
   ```bash
   ./scripts/deploy.sh
   ```

3. **Your API will be live at**: https://nq.nathangrilliot.com

## 📦 What's Included

### Docker Configuration
- **Dockerfile**: Multi-stage build for optimized Go backend (~15MB image)
- **docker-compose.production.yml**: Complete production stack
  - Go backend container with health checks
  - Neo4j database with persistent storage
  - Nginx reverse proxy with SSL termination
  - Certbot for automatic SSL certificate management

### Nginx Configuration
- **SSL/TLS**: Automatic HTTPS with Let's Encrypt certificates
- **Security Headers**: HSTS, CSP, XSS protection, etc.
- **Rate Limiting**: API protection against abuse
- **CORS**: Configured for your domain
- **WebSocket Support**: For GraphQL subscriptions

### Scripts
- **`scripts/deploy.sh`**: Complete production deployment
- **`scripts/update-backend.sh`**: Zero-downtime backend updates
- **`scripts/backup.sh`**: Neo4j database backups
- **`scripts/setup-firewall.sh`**: System security setup
- **`scripts/test-setup.sh`**: Verify configuration before deployment

## 🔐 Security Features

### Network Security
- Docker network isolation
- UFW firewall configuration
- Rate limiting (10 req/s per IP)
- No direct database access from internet

### SSL/TLS
- Automatic certificate generation and renewal
- Strong cipher suites (TLS 1.2+ only)
- HSTS headers with preload ready configuration
- OCSP stapling for improved performance

### Application Security
- JWT authentication (existing)
- CORS protection
- GraphQL query complexity limiting
- Container runs as non-root user
- Security headers on all responses

## 📊 Monitoring & Management

### Health Checks
- **Backend**: `/health` endpoint with Docker health monitoring
- **Neo4j**: Database connectivity monitoring
- **SSL**: Automatic certificate renewal

### Logs
```bash
# View all service logs
docker-compose -f docker-compose.production.yml logs -f

# View specific service logs
docker-compose -f docker-compose.production.yml logs -f backend
docker-compose -f docker-compose.production.yml logs -f nginx
docker-compose -f docker-compose.production.yml logs -f neo4j
```

### Service Management
```bash
# Check status
docker-compose -f docker-compose.production.yml ps

# Restart services
docker-compose -f docker-compose.production.yml restart

# Stop services
docker-compose -f docker-compose.production.yml down

# Update backend only
./scripts/update-backend.sh
```

## 💾 Data Management

### Database Backups
```bash
# Create backup
./scripts/backup.sh

# Backups are stored in ./backups/ with timestamps
# Example: neo4j_backup_20260405_182534.dump
```

### Data Persistence
- **Neo4j data**: Persistent Docker volume (`nq_neo4j_data`)
- **SSL certificates**: Persistent Docker volume (`nq_certbot_certs`)
- **Logs**: Persistent Docker volume (`nq_neo4j_logs`)

### Restore Process
```bash
# Stop Neo4j
docker-compose -f docker-compose.production.yml stop neo4j

# Copy backup into container
docker cp ./backups/neo4j_backup_TIMESTAMP.dump nq_neo4j:/backups/

# Restore database
docker-compose -f docker-compose.production.yml exec neo4j \
    neo4j-admin database load neo4j --from-path=/backups/neo4j_backup_TIMESTAMP.dump --overwrite-destination=true

# Start Neo4j
docker-compose -f docker-compose.production.yml start neo4j
```

## 🔧 Configuration

### Environment Variables (.env.production)
```bash
NEO4J_PASSWORD=hp75gSiVkXTmUCZxYwHtkId8mjbMQ7tlihvqkwtBf3s=
DOMAIN=nq.nathangrilliot.com
COMPOSE_PROJECT_NAME=nq
DOCKER_BUILDKIT=1
```

### Network Configuration
- **HTTP Port**: 80 (redirects to HTTPS)
- **HTTPS Port**: 443 (main API access)
- **Internal Network**: Docker bridge (172.20.0.0/16)

### Resource Allocation
- **Neo4j Heap**: 512m initial, 1G maximum
- **Neo4j Page Cache**: 512m
- **Container Restart Policy**: unless-stopped

## 🚨 Troubleshooting

### SSL Certificate Issues
```bash
# Check certificate status
docker-compose -f docker-compose.production.yml run --rm certbot certificates

# Force certificate renewal
docker-compose -f docker-compose.production.yml run --rm certbot renew --force-renewal

# Restart nginx after certificate changes
docker-compose -f docker-compose.production.yml restart nginx
```

### Container Health Issues
```bash
# Check container health
docker inspect nq_backend --format='{{.State.Health.Status}}'

# View health check logs
docker inspect nq_backend --format='{{range .State.Health.Log}}{{.Output}}{{end}}'

# Manual health check
curl -f http://localhost:8080/health
```

### Database Connection Issues
```bash
# Test Neo4j connectivity
docker-compose -f docker-compose.production.yml exec neo4j cypher-shell -u neo4j -p "$NEO4J_PASSWORD" "RETURN 1"

# Check Neo4j logs
docker-compose -f docker-compose.production.yml logs neo4j
```

### DNS and Connectivity Issues
```bash
# Test DNS resolution
host nq.nathangrilliot.com

# Test external connectivity
curl -I https://nq.nathangrilliot.com/health

# Check firewall rules
sudo ufw status verbose
```

## 🔄 Updates and Maintenance

### Backend Code Updates
```bash
# Pull latest code and update backend
git pull
./scripts/update-backend.sh
```

### System Updates
```bash
# Update system packages
sudo apt update && sudo apt upgrade

# Update Docker images
docker-compose -f docker-compose.production.yml pull
docker-compose -f docker-compose.production.yml up -d
```

### Certificate Renewal
Certificates automatically renew every 12 hours. Manual renewal:
```bash
docker-compose -f docker-compose.production.yml run --rm certbot renew
docker-compose -f docker-compose.production.yml restart nginx
```

## 📞 Support

### API Endpoints
- **Production API**: https://nq.nathangrilliot.com/graphql
- **Health Check**: https://nq.nathangrilliot.com/health
- **GraphQL Playground**: https://nq.nathangrilliot.com/playground

### Useful Resources
- **Neo4j Browser**: Access via port 7474 if needed (restricted to localhost)
- **Docker Logs**: All services log to Docker for centralized viewing
- **System Logs**: Available via journalctl for systemd services

---

**🎉 Your NQ backend is now production-ready and accessible from anywhere!**