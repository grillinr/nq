# Installation Guide

This guide will walk you through setting up NQ on your computer, step-by-step. Even if you're not a developer, you can follow these instructions to get NQ running.

**Time needed:** 30 minutes to 2 hours (depending on experience level)

**What you'll do:**
1. Install required software (Go, Node.js, and a database)
2. Download NQ
3. Set up the backend (the "brain" of NQ)
4. Set up the frontend (the mobile app)
5. Start using NQ!

---

## Step 1: Install Required Software

Before NQ can run, you need to install three programs. Don't worry - they're all free!

### Install Go (Programming Language)

Go powers NQ's backend server.

**What to do:**
1. Visit [https://go.dev/dl/](https://go.dev/dl/)
2. Download the installer for your operating system:
   - Windows: Download the `.msi` file
   - Mac: Download the `.pkg` file
   - Linux: Download the appropriate tarball
3. Run the installer and follow the on-screen instructions
4. Accept the default installation options

**Verify it worked:**
- Open Terminal (Mac/Linux) or Command Prompt (Windows)
- Type: `go version`
- You should see something like: `go version go1.25.1`

**Need help?** See the [detailed Go installation steps in the User Guide](User-Guide.md#step-1-install-go)

### Install Node.js (JavaScript Runtime)

Node.js runs the mobile app development server.

**What to do:**
1. Visit [https://nodejs.org/](https://nodejs.org/)
2. Download the **LTS** (Long Term Support) version - this is the recommended, stable version
3. Run the installer
4. Follow the installation wizard (accept default options)

**Verify it worked:**
- Open Terminal/Command Prompt
- Type: `node --version`
- You should see a version number like: `v18.x.x` or higher

### Set Up Neo4j Database

Neo4j is a special type of database that stores connections between your media (like "this movie is similar to that one").

**Choose one option:**

#### Option A: Neo4j Aura (Cloud) - **Recommended for beginners**

This is easier because nothing installs on your computer.

**What to do:**
1. Go to [https://neo4j.com/cloud/aura/](https://neo4j.com/cloud/aura/)
2. Click "Start Free"
3. Create an account (or sign in with Google)
4. Click "Create Instance" → Select "Free" tier
5. **IMPORTANT**: Save the credentials shown! You'll need:
   - Connection URI (looks like: `neo4j+s://xxxxx.databases.neo4j.io`)
   - Username (usually `neo4j`)
   - Password (auto-generated)
6. Click "Download and Continue" to save these to a text file
7. Wait 1-2 minutes for your database to start

#### Option B: Neo4j Desktop (Local) - **For advanced users**

Install the database on your computer.

**What to do:**
1. Go to [https://neo4j.com/download/](https://neo4j.com/download/)
2. Download Neo4j Desktop for your operating system
3. Install and open Neo4j Desktop
4. Create a new project (name it "NQ")
5. Add a new "Local DBMS" (database)
6. Set a password you'll remember
7. Click "Start" to run your database
8. Your credentials are:
   - Connection URI: `bolt://localhost:7687`
   - Username: `neo4j`
   - Password: [your chosen password]

---

## Step 2: Download NQ

**Option 1: Using Git (Recommended)**

If you have Git installed:
```bash
# Choose where to install (e.g., your home folder)
cd ~

# Download NQ
git clone https://github.com/grillinr/nq.git

# Go into the NQ folder
cd nq
```

**Option 2: Download ZIP**

If you don't have Git:
1. Visit [https://github.com/grillinr/nq](https://github.com/grillinr/nq)
2. Click the green "Code" button
3. Click "Download ZIP"
4. Extract the ZIP file to a location you'll remember
5. Open Terminal/Command Prompt and navigate to that folder

---

## Step 3: Set Up the Backend

The backend is the "server" part of NQ - it talks to your database and the various media services.

### 3.1: Navigate to the Backend Folder

```bash
cd backend
```

(If you're not already in the `nq` folder, navigate there first)

### 3.2: Create Your Configuration File

This file will store your database password and API keys.

**Mac/Linux:**
```bash
cp .envtemplate .env
```

**Windows:**
```cmd
copy .envtemplate .env
```

This creates a new file called `.env` based on the template.

### 3.3: Add Your Database Credentials

**What to do:**
1. Open the `.env` file in a text editor (Notepad on Windows, TextEdit on Mac, or any code editor)
2. Find these lines:
   ```bash
   NEO4J_URI=
   NEO4J_USERNAME=
   NEO4J_PASSWORD=
   NEO4J_DATABASE=
   ```
3. Fill in your Neo4j credentials from Step 1:
   ```bash
   NEO4J_URI=neo4j+s://xxxxx.databases.neo4j.io  # Your connection URI from Aura
   NEO4J_USERNAME=neo4j                           # Usually just "neo4j"
   NEO4J_PASSWORD=YourPasswordHere                # The password you saved
   NEO4J_DATABASE=neo4j                          # Usually just "neo4j"
   ```
4. Save the file

**Note:** Leave the API key sections blank for now. You'll add those later when you want to connect services like Spotify or Steam. See [API Data Sources](API-Data-Sources.md) for details.

### 3.4: Test the Backend

Let's make sure everything works!

```bash
go run .
```

**What should happen:**
- You'll see some text appear
- After a few seconds, you should see: `Starting server on :8080`
- The terminal will stay open and keep running

**Success!** Your backend is working.

**If you see errors:**
- "Failed to connect to Neo4j" → Check your credentials in `.env`
- "Port 8080 already in use" → Close other programs using that port, or see [Troubleshooting](#troubleshooting)
- Other errors → See the [FAQ](FAQ.md) or [Troubleshooting](#troubleshooting)

### 3.5: Verify It's Working

While the backend is running:
1. Open a web browser
2. Go to: `http://localhost:8080`
3. You should see a page called "GraphQL Playground"

**Perfect!** The backend is ready.

**Stop the backend for now:**
- Press `Ctrl+C` in the terminal window

---

## Step 4: Set Up the Frontend

The frontend is the mobile app you'll actually interact with.

### 4.1: Navigate to the Frontend Folder

Open a **new** terminal window (or use the same one now that backend is stopped):

```bash
# From the backend folder:
cd ../frontend

# OR if you're in the main nq folder:
cd frontend
```

### 4.2: Install Dependencies

This downloads all the code libraries the app needs.

```bash
npm install
```

**This will take a few minutes.** You'll see lots of text scroll by - that's normal! It's downloading thousands of small files the app needs.

**Wait until you see:**
- The scrolling stops
- You're back to the command prompt
- No error messages appear

### 4.3: Start the Development Server

```bash
npx expo start
```

**What you'll see:**
- Metro bundler starting (this prepares the app)
- A QR code appears
- A menu with options like "Press a │ open Android"

**This is good!** The frontend is running.

### 4.4: Open the App

You have several ways to view the app:

**Option A: On Your Phone (Easiest)**
1. Install the "Expo Go" app from your phone's app store (it's free)
2. Open Expo Go
3. Scan the QR code shown in your terminal
4. Wait for the app to load (first time takes a minute)

**Option B: Android Emulator**
1. Install [Android Studio](https://developer.android.com/studio)
2. Set up an Android Virtual Device (AVD) - see [Android setup guide](https://docs.expo.dev/workflow/android-studio-emulator/)
3. Start the emulator
4. In the Expo terminal, press `a`

**Option C: iOS Simulator (Mac only)**
1. Install Xcode from the Mac App Store
2. In the Expo terminal, press `i`
3. The simulator will open and load the app

**Option D: Web Browser (Limited)**
1. In the Expo terminal, press `w`
2. A browser window opens (note: some features won't work in browser)

---

## Step 5: Verify Everything Works

### Check the Backend

**While the backend is running** (you started it with `go run .` in Step 3):

1. Open browser to: `http://localhost:8080`
2. You should see "GraphQL Playground"
3. This means the backend is working correctly!

### Check the Frontend

**While the frontend is running** (you started it with `npx expo start` in Step 4):

1. The app should load on your phone/emulator
2. You should see the NQ interface
3. If you see connection errors, make sure the backend is also running

### Test the Full System

To use NQ, you need **both** running at the same time:
- **Terminal 1**: Backend (`cd backend && go run .`)
- **Terminal 2**: Frontend (`cd frontend && npx expo start`)

Keep both terminal windows open while using NQ.

---

## Troubleshooting Common Issues

### "Cannot connect to Neo4j" Error

**Problem:** Backend can't reach your database

**Solutions:**
1. **Check your credentials** in `.env`:
   - Make sure URI, username, and password are exactly as provided
   - No extra spaces or quotes
2. **For Aura (cloud)**: Check that your instance is running at [console.neo4j.io](https://console.neo4j.io)
3. **For local**: Make sure Neo4j Desktop shows your database as "Active"

### "Port 8080 already in use"

**Problem:** Another program is using that port

**Solutions:**
1. **Find and close** the other program using port 8080
2. **OR change NQ's port**: 
   - Open `.env`
   - Add a new line: `PORT=8081`
   - Save the file
   - Restart the backend

### "Command not found: go" or "Command not found: npm"

**Problem:** The software isn't installed or not in your system's PATH

**Solutions:**
1. **Reinstall** Go or Node.js following Step 1
2. **Restart your terminal** after installing
3. **Check the installation** followed all steps

### Frontend Package Errors

**Problem:** Missing or corrupted dependencies

**Solution:**
```bash
cd frontend
rm -rf node_modules  # Delete the folder
npm install           # Reinstall everything
```

### "Metro bundler has encountered an error"

**Problem:** The app bundler has cached corrupted files

**Solution:**
```bash
npx expo start -c  # The -c flag clears the cache
```

### Can't Connect to Backend from Phone

**Problem:** Your phone can't reach localhost

**Solution:**
- Your phone and computer must be on the **same WiFi network**
- Instead of `localhost`, the app may need to use your computer's IP address
- Check your router or network settings to find your computer's local IP (e.g., `192.168.1.100`)

### Still Having Problems?

- **Check the [FAQ](FAQ.md)** for more detailed troubleshooting
- **See the [User Guide](User-Guide.md)** for step-by-step help
- **Review error messages carefully** - they often tell you exactly what's wrong
- **Search for your error online** - many issues have common solutions

---

## Next Steps

🎉 **Congratulations!** NQ is installed and running.

### Now What?

1. **Connect Your Services**
   - See the [API Data Sources](API-Data-Sources.md) guide
   - Get API keys for services you want to use (Spotify, Steam, etc.)
   - Add them to your `.env` file

2. **Learn to Use NQ**
   - Read the [User Guide](User-Guide.md) for detailed instructions
   - Start syncing your media libraries
   - Get personalized recommendations!

3. **Explore the Features**
   - Try the GraphQL Playground at `http://localhost:8080`
   - Browse your media library in the app
   - Discover new content based on your preferences

### Daily Usage

Whenever you want to use NQ:

1. **Start the backend:**
   ```bash
   cd /path/to/nq/backend
   go run .
   ```

2. **Start the frontend** (in a new terminal):
   ```bash
   cd /path/to/nq/frontend
   npx expo start
   ```

3. **Open the app** on your phone or emulator

4. **Enjoy your recommendations!**

### Keep Learning

- **[User Guide](User-Guide.md)** - Complete how-to for using NQ
- **[FAQ](FAQ.md)** - Answers to common questions
- **[API Data Sources](API-Data-Sources.md)** - Connect your media services
- **[Usage Guide](Usage.md)** - Advanced usage and development
- **[Tech Stack](Tech-Stack.md)** - Learn about the technology behind NQ
