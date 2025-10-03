# ✅ Universal Deployment Checklist

Use this checklist for deploying MH-Go to **any platform** that supports Docker.

---

## 📋 Pre-Deployment Checklist

### 1. Code & Repository

- [ ] All changes committed to Git
- [ ] Code pushed to GitHub/GitLab/Bitbucket
- [ ] `.env` file is in `.gitignore` (never commit secrets!)
- [ ] Tests passing locally

### 2. Environment Configuration

- [ ] Generate JWT secret: `openssl rand -base64 32`
- [ ] Prepare database connection string
- [ ] (Optional) Configure SMTP for email features
- [ ] (Optional) Get OpenAI API key for AI features

### 3. Platform Selection

Choose your deployment platform:

- [ ] Cloud VM (AWS, GCP, Azure, DigitalOcean)
- [ ] Container Service (ECS, Cloud Run, ACI)
- [ ] Kubernetes cluster
- [ ] PaaS (Railway, Render, Fly.io, Heroku)
- [ ] Self-hosted server

---

## 🚀 Deployment Steps (Universal)

### Step 1: Set Up Database

**Option A: Platform's Managed PostgreSQL**

- [ ] Create PostgreSQL instance on your platform
- [ ] Note the `DATABASE_URL` connection string
- [ ] Ensure SSL is enabled (`?sslmode=require`)

**Option B: External PostgreSQL Service**

- [ ] Sign up for: AWS RDS, Supabase, Neon, ElephantSQL, etc.
- [ ] Create database instance
- [ ] Copy connection string

**Option C: Self-Hosted PostgreSQL**

- [ ] Use included `docker-compose.yml`
- [ ] Run: `docker-compose up -d postgres`
- [ ] Connection: `postgres://user:password@localhost:5432/mhgo`

### Step 2: Configure Environment Variables

Set these variables in your platform (exact location varies):

**Required:**

```bash
PORT=8080
ENVIRONMENT=production
JWT_SECRET=<your-generated-secret>
DATABASE_URL=<your-postgres-connection-string>
```

**Optional:**

```bash
# For email features
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
FROM_EMAIL=noreply@yourdomain.com

# For AI features
OPENAI_API_KEY=sk-your-key-here
```

### Step 3: Deploy Application

**For PaaS (Railway, Render, Fly.io, Heroku):**

- [ ] Connect GitHub repository
- [ ] Platform auto-detects Dockerfile
- [ ] Click "Deploy"
- [ ] Wait for build to complete

**For Cloud VMs / Self-Hosted:**

```bash
# 1. SSH into server
ssh user@your-server-ip

# 2. Install Docker (if not installed)
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 3. Clone repository
git clone https://github.com/your-username/mh-go.git
cd mh-go

# 4. Configure environment
cp .env.example .env
nano .env  # Edit with your values

# 5. Deploy
docker-compose -f docker-compose.prod.yml up -d
```

**For Kubernetes:**

- [ ] Build and push image to registry
- [ ] Create secrets: `kubectl create secret generic mhgo-secrets`
- [ ] Apply deployment: `kubectl apply -f k8s/`
- [ ] Verify pods: `kubectl get pods`

**For Container Services (ECS, Cloud Run, ACI):**

- [ ] Build Docker image locally or in CI
- [ ] Push to container registry (ECR, GCR, ACR)
- [ ] Create service with environment variables
- [ ] Configure health checks (`/health`)
- [ ] Deploy

### Step 4: Verify Deployment

- [ ] Check build logs for errors
- [ ] Verify containers/pods are running
- [ ] Test health endpoint: `curl https://your-domain.com/health`
  - Expected: `{"status":"ok","id":"..."}`
- [ ] Check application logs for "Connected to database successfully"
- [ ] Open application in browser
- [ ] Test registration and login

### Step 5: Post-Deployment

- [ ] Set up custom domain (if using one)
- [ ] Configure SSL/HTTPS
- [ ] Set up monitoring and alerts
- [ ] Configure automated backups
- [ ] Document your deployment process
- [ ] Test all critical features

---

## 🔍 Platform-Specific Quick Links

### PaaS Platforms

**Railway:**

1. Connect GitHub → Select repo
2. Add PostgreSQL database
3. Set env vars: `DATABASE_URL=${{Postgres.DATABASE_URL}}`
4. Generate domain

**Render:**

1. New Web Service → Connect repo
2. Add PostgreSQL database
3. Set environment variables
4. Deploy

**Fly.io:**

```bash
flyctl launch
flyctl postgres create
flyctl postgres attach <db-name>
flyctl deploy
```

**Heroku:**

```bash
heroku create your-app-name
heroku addons:create heroku-postgresql:mini
heroku config:set JWT_SECRET=$(openssl rand -base64 32)
git push heroku main
```

### Cloud Providers

**AWS ECS:**

1. Push image to ECR
2. Create task definition
3. Create service with ALB
4. Use RDS for PostgreSQL

**Google Cloud Run:**

```bash
gcloud builds submit --tag gcr.io/PROJECT/mh-go
gcloud run deploy mh-go --image gcr.io/PROJECT/mh-go
```

**Azure Container Instances:**

```bash
az container create --resource-group myRG \
  --name mh-go --image your-registry/mh-go:latest \
  --environment-variables PORT=8080 DATABASE_URL=...
```

**DigitalOcean App Platform:**

1. Connect GitHub repo
2. Auto-detects Dockerfile
3. Add managed PostgreSQL
4. Set environment variables
5. Deploy

---

## 🐛 Troubleshooting

### Build Fails

**Check:**

- [ ] Dockerfile syntax
- [ ] All dependencies in `go.mod`
- [ ] Test build locally: `docker build -t mh-go-test .`

**Common issues:**

- Missing templ CLI installation in Dockerfile ✅ (already fixed)
- Network issues downloading dependencies
- Wrong Go version specified

### Application Won't Start

**Check logs for:**

- `DATABASE_URL` not set or incorrect format
- `JWT_SECRET` not set
- Port already in use
- Database not reachable

**Solutions:**

- Verify all environment variables are set
- Check database connection string format
- Ensure database is running and accessible
- Check firewall/security group rules

### Database Connection Errors

**Error: "no such host"**

- Using wrong hostname (e.g., `localhost` in container)
- Use platform's provided `DATABASE_URL`
- Or use container name in Docker network

**Error: "connection refused"**

- Database not running
- Wrong port
- Firewall blocking connection
- Need to add `?sslmode=require` or `?sslmode=disable`

**Error: "authentication failed"**

- Wrong username/password
- Check credentials in `DATABASE_URL`

### Application Crashes

**Check:**

- [ ] View application logs
- [ ] Check database migrations status
- [ ] Verify all required env vars are set
- [ ] Test health endpoint: `/health`

---

## 📊 Success Criteria

Your deployment is successful when:

✅ Build completes without errors
✅ Application starts successfully
✅ Health check returns `{"status":"ok"}`
✅ Database migrations run automatically
✅ Can access application in browser
✅ Can register and login
✅ Can create folders and content
✅ Application auto-restarts on failure

---

## 📚 Additional Resources

- **Quick Start**: [DOCKER_README.md](./DOCKER_README.md)
- **Full Guide**: [DEPLOYMENT.md](./DEPLOYMENT.md)
- **Environment Template**: [.env.example](./.env.example)
- **Docker Compose**: [docker-compose.yml](./docker-compose.yml)
- **Production Config**: [docker-compose.prod.yml](./docker-compose.prod.yml)

---

## 🆘 Getting Help

If you encounter issues:

1. **Check logs** (most informative source)
2. **Test locally first** with Docker Compose
3. **Verify environment variables** are set correctly
4. **Check database connectivity** separately
5. **Review platform-specific documentation**
6. **Create GitHub issue** with logs and details

---

**Remember:** This application uses standard Docker. If it works locally with Docker, it will work anywhere! 🚀
