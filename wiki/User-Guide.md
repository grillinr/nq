# User Guide

Welcome to the NQ User Guide! This guide will walk you through everything you need to know to set up and use NQ, from installation to getting your first recommendations.

**What you'll learn:**
- How to install NQ from scratch
- How to connect your media accounts
- How to get recommendations
- How to troubleshoot common issues

**Before you start:**
- This guide assumes you're comfortable using a computer and following step-by-step instructions
- You'll need to use the command line (Terminal on Mac/Linux, Command Prompt or PowerShell on Windows)
- Budget about 1-2 hours for first-time setup

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Installing Prerequisites](#installing-prerequisites)
3. [Setting Up Your Database](#setting-up-your-database)
4. [Installing NQ](#installing-nq)
5. [Connecting Your Services](#connecting-your-services)
6. [Running NQ](#running-nq)
7. [Using NQ](#using-nq)
8. [Troubleshooting](#troubleshooting)

---

## Getting Started

### What is NQ?

NQ (Next in Queue) helps you decide what to watch, play, read, or listen to next. Instead of scrolling endlessly through Netflix or your game library, NQ analyzes your preferences and suggests content you'll actually enjoy.

### How does it work?

1. **Connect your accounts**: Link services like Spotify, Steam, YouTube
2. **Sync your data**: NQ downloads your watch/play/listen history
3. **Build your graph**: Creates a map of your preferences and relationships
4. **Get recommendations**: Receive personalized suggestions for what's next

### What you'll need

- A computer (Windows, Mac, or Linux)
- Internet connection
- Accounts on services you want to connect (Spotify, Steam, etc.)
- About 2GB of free disk space

---

## Installing Prerequisites

NQ requires three software programs before it can run. Follow the instructions for your operating system.

### Step 1: Install Go

Go is the programming language NQ's backend uses.

**Windows:**
1. Visit https://go.dev/dl/
2. Download the Windows installer (`.msi` file)
3. Run the installer and follow the prompts
4. Accept default settings

**Mac:**
1. Visit https://go.dev/dl/
2. Download the macOS installer (`.pkg` file)
3. Double-click the file and follow installation prompts
   
   *OR* use Homebrew:
   ```bash
   brew install go
   ```

**Linux:**
1. Visit https://go.dev/dl/
2. Download the Linux tarball
3. Extract and install:
   ```bash
   sudo tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   ```

**Verify installation:**
```bash
go version
```
You should see: `go version go1.25.1` (or similar)

### Step 2: Install Node.js

Node.js runs the mobile app development server.

**Windows:**
1. Visit https://nodejs.org/
2. Download the LTS (Long Term Support) version
3. Run the installer
4. Check "Automatically install necessary tools" during setup

**Mac:**
1. Visit https://nodejs.org/
2. Download the LTS version installer
3. Run the `.pkg` file
   
   *OR* use Homebrew:
   ```bash
   brew install node
   ```

**Linux:**
```bash
# Ubuntu/Debian
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs

# Fedora
sudo dnf install nodejs

# Arch
sudo pacman -S nodejs npm
```

**Verify installation:**
```bash
node --version
npm --version
```
You should see version numbers (18 or higher for Node.js)

### Step 3: Install Git (Optional but Recommended)

Git helps you download and update NQ.

**Windows:**
1. Visit https://git-scm.com/download/win
2. Download and run the installer
3. Use default settings

**Mac:**
1. Open Terminal
2. Type: `git --version`
3. If not installed, Mac will prompt you to install Command Line Tools

**Linux:**
```bash
# Ubuntu/Debian
sudo apt-get install git

# Fedora
sudo dnf install git

# Arch
sudo pacman -S git
```

**Verify installation:**
```bash
git --version
```

---

## Setting Up Your Database

NQ uses Neo4j to store your media data and relationships. You have two options: cloud or local.

### Option 1: Neo4j Aura (Cloud) - Recommended for Beginners

**Advantages:**
- Free tier available
- No software to install
- Automatic backups
- Works from anywhere

**Setup Steps:**

1. **Create an account**
   - Visit https://neo4j.com/cloud/aura/
   - Click "Start Free"
   - Sign up with email or Google account

2. **Create a free instance**
   - Click "Create Instance"
   - Choose "Free" tier
   - Select a region close to you
   - Click "Create"

3. **Save your credentials**
   - **IMPORTANT**: A dialog will appear with your credentials
   - Copy these somewhere safe:
     - Connection URI (looks like: `neo4j+s://xxxxx.databases.neo4j.io`)
     - Username (usually `neo4j`)
     - Generated password
   - You **cannot** retrieve the password later!
   - Click "Download and Continue" to save the `.txt` file

4. **Wait for instance to start**
   - Status will show "Running" when ready (takes 1-2 minutes)

### Option 2: Neo4j Desktop (Local) - For Advanced Users

**Advantages:**
- Full control over data
- Works offline
- No account needed

**Setup Steps:**

1. **Download Neo4j Desktop**
   - Visit https://neo4j.com/download/
   - Download Neo4j Desktop for your OS
   - Install the application

2. **Create a database**
   - Open Neo4j Desktop
   - Click "New" → "Create Project"
   - Name your project (e.g., "NQ")
   - Click "Add" → "Local DBMS"
   - Name it (e.g., "NQ Database")
   - Set a password (remember this!)
   - Click "Create"

3. **Start the database**
   - Click "Start" on your database
   - Wait for it to show "Active"

4. **Note your credentials**
   - Connection URI: `bolt://localhost:7687`
   - Username: `neo4j`
   - Password: [your chosen password]

---

## Installing NQ

### Step 1: Download NQ

**Using Git (Recommended):**
```bash
# Navigate to where you want NQ installed
cd ~
# or on Windows: cd C:\Users\YourName\

# Clone the repository
git clone https://github.com/grillinr/nq.git

# Navigate into the folder
cd nq
```

**Without Git:**
1. Visit https://github.com/grillinr/nq
2. Click the green "Code" button
3. Click "Download ZIP"
4. Extract the ZIP file
5. Open Terminal/Command Prompt and navigate to the extracted folder

### Step 2: Set Up the Backend

1. **Navigate to the backend folder**
   ```bash
   cd backend
   ```

2. **Create your configuration file**
   
   **Mac/Linux:**
   ```bash
   cp .envtemplate .env
   ```
   
   **Windows:**
   ```cmd
   copy .envtemplate .env
   ```

3. **Edit the .env file**
   
   Open `.env` in a text editor (Notepad, TextEdit, VS Code, etc.)
   
   **Required settings** - Fill these in now:
   ```bash
   # Neo4j Database Connection
   NEO4J_URI=neo4j+s://xxxxx.databases.neo4j.io  # Your connection URI
   NEO4J_USERNAME=neo4j                           # Usually 'neo4j'
   NEO4J_PASSWORD=your_password_here             # Your database password
   NEO4J_DATABASE=neo4j                          # Usually 'neo4j'
   ```
   
   **Optional API settings** - Add these later when you connect services:
   ```bash
   # Example: Spotify (add when you're ready)
   SPOTIFY_CLIENT_ID=
   SPOTIFY_CLIENT_SECRET=
   
   # Example: Steam (add when you're ready)
   STEAM_API_KEY=
   
   # Leave blank any services you don't plan to use
   ```

4. **Test the backend**
   ```bash
   go run .
   ```
   
   **What you should see:**
   ```
   Starting server on :8080
   ```
   
   **If you see errors:**
   - Check your Neo4j credentials in `.env`
   - Make sure your Neo4j instance is running
   - See [Troubleshooting](#troubleshooting)

5. **Stop the backend** (for now)
   - Press `Ctrl+C` in the terminal

### Step 3: Set Up the Frontend

1. **Navigate to the frontend folder**
   ```bash
   # From the backend folder:
   cd ../nq-frontend
   
   # Or from the nq root folder:
   cd nq-frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```
   
   This will take a few minutes and download many files. This is normal!

3. **Test the frontend**
   ```bash
   npx expo start
   ```
   
   **What you should see:**
   - Metro bundler starting
   - A QR code
   - Instructions for opening the app
   
   **Success!** The frontend is working.

4. **Stop the frontend** (for now)
   - Press `Ctrl+C` in the terminal

---

## Connecting Your Services

Now the fun part - connecting your media services! You only need to connect services you actually use.

### Understanding API Keys

Think of API keys like special passwords that let NQ access your data on other services. Each service has its own process for getting these keys.

**Important:**
- Most APIs are free for personal use
- You're only accessing your own data
- Keys are stored securely in your `.env` file
- You can add services gradually

### Step-by-Step: Connecting Spotify

Let's walk through connecting Spotify as an example.

1. **Visit Spotify for Developers**
   - Go to https://developer.spotify.com/dashboard
   - Log in with your Spotify account

2. **Create an app**
   - Click "Create app"
   - Fill in the form:
     - App name: `NQ` (or anything you like)
     - App description: `Personal media recommendation system`
     - Redirect URI: `http://localhost:8080` (required but not used yet)
     - Which API are you using?: `Web API`
   - Check the terms of service box
   - Click "Save"

3. **Get your credentials**
   - Click on your new app
   - Click "Settings"
   - You'll see:
     - **Client ID**: A long string of letters and numbers
     - **Client Secret**: Click "View client secret" to reveal it
   - Copy both of these

4. **Add to your .env file**
   - Open `backend/.env` in a text editor
   - Find the Spotify section:
     ```bash
     SPOTIFY_CLIENT_ID=paste_client_id_here
     SPOTIFY_CLIENT_SECRET=paste_client_secret_here
     ```
   - Save the file

5. **Restart your backend**
   ```bash
   cd backend
   go run .
   ```
   
   NQ will now be able to sync your Spotify data!

### Connecting Other Services

For detailed instructions on other services, see the [API Data Sources](API-Data-Sources.md) guide.

**Quick links to API setup pages:**
- **Steam**: https://steamcommunity.com/dev/apikey
- **YouTube**: https://console.cloud.google.com/ (Enable YouTube Data API v3)
- **Twitch**: https://dev.twitch.tv/console
- **TMDB**: https://www.themoviedb.org/settings/api
- **Instapaper**: Just use your regular login credentials
- **Open Library**: No credentials needed!

**You don't need to connect everything!** Start with 1-2 services and add more later.

---

## Running NQ

### Starting the Backend

The backend is the "brain" of NQ - it manages your data and provides recommendations.

1. **Open a terminal**
2. **Navigate to the backend folder**
   ```bash
   cd /path/to/nq/backend
   ```
3. **Start the server**
   ```bash
   go run .
   ```
4. **Leave this terminal window open** - the backend needs to keep running

**You'll know it's working when you see:**
```
Starting server on :8080
```

**Test it:**
- Open a web browser
- Go to http://localhost:8080
- You should see the GraphQL Playground

### Starting the Frontend

The frontend is the mobile app you'll interact with.

1. **Open a NEW terminal window** (keep the backend running in the other one)
2. **Navigate to the frontend folder**
   ```bash
   cd /path/to/nq/nq-frontend
   ```
3. **Start Expo**
   ```bash
   npx expo start
   ```

**You'll see an interactive menu:**
```
› Press a │ open Android
› Press i │ open iOS simulator
› Press w │ open web

› Press r │ reload app
› Press m │ toggle menu
```

### Opening the App

You have several options:

**Option 1: On your phone (Easiest)**
1. Install "Expo Go" app from your app store
2. Open Expo Go
3. Scan the QR code shown in your terminal
4. The app will load on your phone

**Option 2: Android Emulator**
1. Install Android Studio
2. Set up an Android Virtual Device (AVD)
3. Press `a` in the Expo terminal
4. The app will open in the emulator

**Option 3: iOS Simulator (Mac only)**
1. Install Xcode from the Mac App Store
2. Press `i` in the Expo terminal
3. The app will open in the simulator

**Option 4: Web Browser**
1. Press `w` in the Expo terminal
2. The app opens in your browser (limited functionality)

---

## Using NQ

### First Time Setup

When you first open the app:

1. **Login/Signup** (if authentication is enabled)
   - Create an account or log in

2. **Connect Services**
   - Navigate to Settings → Integrations
   - You'll see a list of available services
   - Services with API keys in your `.env` will be enabled

3. **Sync Your Data**
   - Tap "Sync" for each service you want to use
   - First sync takes a while (5-30 minutes depending on library size)
   - Progress indicators show the sync status

4. **Wait for Completion**
   - You'll get a notification when syncing is complete
   - You can use NQ while syncing continues in the background

### Getting Recommendations

1. **Open the Home/Queue Screen**
   - This is the main screen of the app

2. **View Your Queue**
   - NQ displays personalized recommendations
   - Items are ranked by "fitness score" (how well they match your preferences)

3. **Explore Recommendations**
   - Tap on an item to see details
   - See why it was recommended (similar to what you've liked)
   - View trailers, descriptions, ratings, etc.

4. **Take Action**
   - Mark as "Want to Watch/Play/Read"
   - Mark as "Not Interested"
   - Add to custom lists

### Browsing Your Library

1. **Navigate to Library**
   - See all media items from your connected services
   - Filter by type (movies, games, music, books)
   - Sort by various criteria

2. **View Item Details**
   - Tap any item to see full details
   - View cast, creators, release date, etc.
   - See related items and recommendations

### Managing Integrations

1. **Go to Settings → Integrations**

2. **Sync Individual Services**
   - Tap "Sync" next to any service
   - Useful for updating specific data

3. **Disconnect Services**
   - Remove services you no longer want to track
   - Data remains in your database unless you delete it

### Using the GraphQL Playground (Advanced)

For direct API access:

1. **Make sure the backend is running**

2. **Open http://localhost:8080 in your browser**

3. **Try a query:**
   ```graphql
   query {
     users {
       id
       username
       email
     }
   }
   ```

4. **Click the Play button**
   - Results appear on the right side

5. **Explore the schema**
   - Click "Docs" on the right to see all available queries
   - Auto-complete helps you write queries (Ctrl+Space)

---

## Troubleshooting

### Backend Won't Start

**"Port 8080 already in use"**
- Another program is using port 8080
- **Solution**: Stop the other program, or change NQ's port:
  - Edit `.env` and add: `PORT=8081`
  - Update your frontend to connect to port 8081

**"Failed to connect to Neo4j"**
- Your database credentials are incorrect or database is offline
- **Solution**: 
  - Verify credentials in `.env` match your Neo4j instance
  - For Aura: Check that your instance is running in the Aura console
  - For local: Open Neo4j Desktop and start your database

**"No such file or directory: .env"**
- You forgot to create the .env file
- **Solution**: `cp .envtemplate .env` (or `copy` on Windows)

### Frontend Won't Start

**"Cannot find module" errors**
- Dependencies aren't installed
- **Solution**: `npm install` in the nq-frontend folder

**"Metro bundler exited"**
- Cache corruption or dependency conflict
- **Solution**: 
  ```bash
  npx expo start -c   # Clear cache
  # If that doesn't work:
  rm -rf node_modules
  npm install
  ```

**"Unable to resolve module"**
- Usually a caching issue
- **Solution**: `npx expo start -c`

### Can't Connect to Backend from App

**"Network request failed"**
- Backend isn't running or wrong address
- **Solution**:
  - Make sure backend is running: `cd backend && go run .`
  - Check the app's API endpoint configuration
  - For physical devices: Use your computer's IP address instead of localhost

**On Android/iOS Emulator:**
- Use `http://10.0.2.2:8080` (Android) or `http://localhost:8080` (iOS)

**On Physical Device:**
- Use your computer's local IP (e.g., `http://192.168.1.100:8080`)
- Make sure your phone and computer are on the same WiFi network

### Sync Issues

**"Sync failed" errors**
- API credentials might be wrong or expired
- **Solution**:
  - Verify API keys in `.env` are correct
  - Check for extra spaces or quotes
  - Regenerate the API key and try again

**Sync never completes**
- Large libraries take time; rate limiting slows things down
- **Solution**:
  - Be patient - first sync can take 30+ minutes
  - Check backend logs for errors
  - Try syncing one service at a time

**"Rate limit exceeded"**
- You've made too many API requests
- **Solution**:
  - Wait 15-60 minutes and try again
  - NQ implements automatic retry with backoff

### Data Issues

**"No recommendations showing"**
- Need more data synced first
- **Solution**:
  - Sync at least one service completely
  - Make sure you have consumed media (played games, listened to music, etc.)
  - Check that metadata enrichment completed

**Missing media items**
- Not all items have metadata available
- **Solution**: 
  - Some items might not be in metadata databases
  - Check backend logs to see if items were skipped

---

## Next Steps

### Learn More

- **[FAQ](FAQ.md)**: Answers to common questions
- **[API Data Sources](API-Data-Sources.md)**: Detailed API setup for all services
- **[Usage Guide](Usage.md)**: Development workflow and advanced usage
- **[Tech Stack](Tech-Stack.md)**: Learn about the technologies NQ uses

### Get Involved

- **Report bugs**: Submit issues on GitHub
- **Request features**: Open a feature request issue
- **Contribute**: Pull requests welcome!

### Optimize Your Experience

1. **Connect more services** for better recommendations
2. **Regularly sync** to keep data up to date
3. **Provide feedback** on recommendations to improve accuracy
4. **Explore the graph** to discover interesting connections

---

## Tips and Best Practices

### For Best Results

- **Connect multiple services in the same media category** (e.g., both Spotify and Apple Music) for richer music data
- **Sync regularly** - weekly or after major consumption sessions
- **Interact with recommendations** - marking items helps NQ learn your preferences
- **Check the GraphQL playground** to explore your data and relationships

### Performance Tips

- **Sync during off-hours** if you have large libraries (overnight)
- **Sync one service at a time** if experiencing issues
- **Use a local Neo4j instance** for faster queries if you're technical
- **Keep services you don't use disconnected** to reduce sync time

### Privacy Tips

- **Never commit your .env file** to version control
- **Use strong passwords** for your Neo4j database
- **Regularly review connected services** and disconnect unused ones
- **Back up your Neo4j database** to prevent data loss

---

**Congratulations!** You're now ready to use NQ. Enjoy discovering your next favorite media!

For questions not covered here, check the [FAQ](FAQ.md) or submit an issue on GitHub.
