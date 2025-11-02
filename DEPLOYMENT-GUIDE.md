# 📂 Complete Netlify V2Ray Proxy Project

## Folder Structure

```
/project-root (dnstt repository)
├── netlify.toml                  # Netlify configuration
├── package.json                  # Node.js package metadata
├── functions/
│   └── proxy.js                  # Netlify Function for reverse proxy
├── NETLIFY-PROXY-README.md       # Complete documentation
└── DEPLOYMENT-GUIDE.md           # This file
```

## 🧾 File Summaries

### 1. netlify.toml
- Defines the functions directory
- Routes `/xhttp/*` requests to the proxy function
- Sets up CORS headers for all methods
- Uses esbuild for function bundling

### 2. package.json
- Node.js project metadata
- No external dependencies (uses native fetch API)
- Requires Node.js 18.0.0 or higher

### 3. functions/proxy.js
- Netlify serverless function
- Forwards requests to `https://ad.sdupdates.news`
- Preserves HTTP method, headers, body, and query parameters
- Returns responses with proper status codes
- Handles errors gracefully

## 🚀 Quick Start

### Deploy to Netlify (3 steps)

1. **Install Netlify CLI**:
   ```bash
   npm install -g netlify-cli
   ```

2. **Login and Initialize**:
   ```bash
   netlify login
   netlify init
   ```

3. **Deploy**:
   ```bash
   netlify deploy --prod
   ```

Your proxy will be live at: `https://your-site-name.netlify.app/xhttp`

## 🔍 How It Works

**Request Flow:**
```
V2Ray Client
    ↓ (HTTPS)
Netlify CDN (https://mydomain.netlify.app/xhttp)
    ↓ (Netlify Function)
Reverse Proxy Function (functions/proxy.js)
    ↓ (HTTPS)
Origin Server (https://ad.sdupdates.news/xhttp)
    ↓ (Response)
Back to V2Ray Client
```

**Key Features:**
- ✅ Hides origin server IP
- ✅ Full HTTPS/TLS encryption
- ✅ Preserves all headers and body
- ✅ Supports all HTTP methods
- ✅ CDN edge caching where appropriate

## ⚙️ V2Ray Client Setup

Update your V2Ray client with:

```json
{
  "v": "2",
  "ps": "Netlify-CDN-V2Ray",
  "add": "your-site-name.netlify.app",
  "port": "443",
  "id": "your-uuid-here",
  "aid": "0",
  "scy": "none",
  "net": "http",
  "type": "none",
  "host": "your-site-name.netlify.app",
  "path": "/xhttp",
  "tls": "tls",
  "sni": "your-site-name.netlify.app"
}
```

**Replace:**
- `your-site-name` with your actual Netlify site name
- `your-uuid-here` with your V2Ray UUID

## 🧪 Testing

```bash
# Test basic connectivity
curl -i https://your-site-name.netlify.app/xhttp -v

# Test with custom headers
curl https://your-site-name.netlify.app/xhttp \
  -H "User-Agent: V2Ray/4.0" \
  -H "Accept: */*"

# Test POST request
curl https://your-site-name.netlify.app/xhttp \
  -X POST \
  -H "Content-Type: application/octet-stream" \
  -d "test data"
```

## ⚠️ Important Limitations

- **Timeout**: 10 seconds per request (free plan)
- **Body Size**: 10 MB maximum
- **Invocations**: 125K per month (free plan)
- **Not suitable for**: Long-lived connections, WebSockets, or streaming

For complete documentation, see **NETLIFY-PROXY-README.md**

## 📚 Additional Files

- **NETLIFY-PROXY-README.md** - Complete technical documentation with all details
- **netlify.toml** - Configuration file
- **package.json** - Package metadata
- **functions/proxy.js** - Proxy implementation

---

**✅ Ready to Deploy!** All files are in place and ready for deployment.
