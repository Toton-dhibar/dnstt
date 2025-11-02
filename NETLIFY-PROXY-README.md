# Netlify Reverse Proxy for V2Ray (xhttp + TLS)

This project provides a complete, working **Netlify CDN reverse proxy** setup for V2Ray xhttp + TLS traffic. It allows you to route your V2Ray traffic through Netlify CDN, hiding your real server IP address while maintaining full HTTPS security.

## 📂 Folder Structure

```
/project-root
├── netlify.toml          # Netlify configuration
├── package.json          # Node.js dependencies and metadata
├── functions/
│   └── proxy.js         # Netlify Function for proxying requests
└── NETLIFY-PROXY-README.md  # This documentation
```

## 🎯 Current Setup

- **Origin V2Ray Server**: `https://ad.sdupdates.news`
- **V2Ray Protocol**: vless + xhttp + TLS
- **Path**: `/xhttp`
- **Port**: 443
- **Netlify Endpoint**: `https://mydomain.netlify.app/xhttp` (forwards to origin)

## 🔍 How It Works

### Reverse Proxy Architecture

The Netlify Function acts as a **middle layer** that forwards client requests to the origin server and returns the response:

1. **Client Request** → Sends request to `https://mydomain.netlify.app/xhttp`
2. **Netlify CDN** → Handles TLS termination and routes to Netlify Function
3. **Netlify Function** → Forwards request to `https://ad.sdupdates.news/xhttp`
4. **Origin Server** → Processes V2Ray request and returns response
5. **Netlify Function** → Returns response to client

### Connection Flow

```
Client → Netlify CDN → Netlify Function → Origin (ad.sdupdates.news)
         [TLS]           [Proxy]            [TLS]
```

### Why This Hides the Real IP

- The end client **only communicates with the Netlify domain**
- Netlify makes the outgoing request to the origin **privately from Netlify's infrastructure**
- Your origin server IP remains hidden from the client
- Only Netlify's IP addresses are visible to the origin (can be whitelisted if needed)

## 🚀 Deployment Instructions

### Method 1: Deploy via Netlify CLI

1. **Install Netlify CLI** (if not already installed):
   ```bash
   npm install -g netlify-cli
   ```

2. **Login to Netlify**:
   ```bash
   netlify login
   ```

3. **Initialize the project** (first time only):
   ```bash
   netlify init
   ```
   - Choose "Create & configure a new site"
   - Select your team
   - Choose a site name (e.g., `mydomain`)

4. **Deploy to production**:
   ```bash
   netlify deploy --prod
   ```

5. **Your site will be available at**: `https://mydomain.netlify.app`

### Method 2: Deploy via GitHub + Netlify Dashboard

1. **Push this repository to GitHub**

2. **Login to Netlify Dashboard** at https://app.netlify.com

3. **Click "Add new site" → "Import an existing project"**

4. **Connect your GitHub repository**

5. **Configure build settings**:
   - Build command: *(leave empty)*
   - Publish directory: *(leave empty)*
   - Functions directory: `functions` (should auto-detect)

6. **Deploy**

7. **Automatic deployments** will trigger on every push to your GitHub repository

### Method 3: Deploy via Netlify Drop

1. **Zip the entire project folder** (including `netlify.toml`, `package.json`, and `functions/`)

2. **Visit** https://app.netlify.com/drop

3. **Drag and drop** your zip file

4. **Wait for deployment** to complete

## 🧪 Testing the Endpoint

After deployment, test your proxy endpoint:

```bash
# Basic connectivity test
curl -i https://mydomain.netlify.app/xhttp/test -v

# Test with headers
curl -i https://mydomain.netlify.app/xhttp \
  -H "User-Agent: V2Ray/4.0" \
  -H "Accept: */*" \
  -v

# POST request test
curl -i https://mydomain.netlify.app/xhttp \
  -X POST \
  -H "Content-Type: application/octet-stream" \
  -d "test data" \
  -v
```

**Expected Behavior**: The proxy should forward your request to `https://ad.sdupdates.news/xhttp/test` and return the origin's response.

## ⚙️ V2Ray Client Configuration

Configure your V2Ray client to use the Netlify proxy domain instead of connecting directly to the origin:

### JSON Configuration

```json
{
  "v": "2",
  "ps": "Netlify-CDN-V2Ray",
  "add": "mydomain.netlify.app",
  "port": "443",
  "id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "aid": "0",
  "scy": "none",
  "net": "http",
  "type": "none",
  "host": "mydomain.netlify.app",
  "path": "/xhttp",
  "tls": "tls",
  "sni": "mydomain.netlify.app"
}
```

### Configuration Parameters Explained

- **add**: Your Netlify domain (e.g., `mydomain.netlify.app`)
- **port**: 443 (HTTPS)
- **id**: Your V2Ray UUID (same as your origin server)
- **net**: `http` (for xhttp protocol)
- **path**: `/xhttp` (must match your V2Ray server path)
- **tls**: `tls` (enable TLS encryption)
- **host**: Your Netlify domain
- **sni**: Server Name Indication (your Netlify domain)

### V2Ray Client App Configuration

For V2Ray GUI clients (V2RayN, V2RayNG, etc.):

1. **Server Address**: `mydomain.netlify.app`
2. **Port**: `443`
3. **User ID**: *(your V2Ray UUID)*
4. **Protocol**: VLESS
5. **Network**: HTTP (xhttp)
6. **Path**: `/xhttp`
7. **TLS**: Enabled
8. **SNI**: `mydomain.netlify.app`

## ⚠️ Netlify Function Limits

Be aware of these limitations on Netlify's free plan:

| Limit | Free Plan | Pro Plan |
|-------|-----------|----------|
| **Max execution time** | 10 seconds | 26 seconds |
| **Request body size** | 10 MB | 10 MB |
| **Response size** | 10 MB | 10 MB |
| **Function invocations** | 125K/month | 1M/month |
| **Bandwidth** | 100 GB/month | 1 TB/month |

### Important Considerations

- ⏱️ **Timeout**: Functions timeout after 10 seconds (free) or 26 seconds (pro). This is suitable for most V2Ray connections, but very long requests may fail.

- 📦 **Size Limits**: Request and response bodies are limited to 10 MB. Large file transfers may not work.

- 🔌 **Idle Connections**: Netlify Functions don't support long-lived TCP connections or WebSocket upgrades natively.

- 🚫 **Not Suitable For**:
  - Long-lived streaming connections
  - WebSocket connections
  - Very large file transfers (>10 MB)
  - High-frequency polling

### Alternative Solutions for Limitations

If you need features beyond Netlify's limits, consider these alternatives:

- **Cloudflare Workers**: 50ms CPU time, better for streaming
- **Fly.io**: Full TCP support, persistent connections
- **AWS Lambda + API Gateway**: Configurable timeouts up to 15 minutes
- **Self-hosted reverse proxy**: Nginx, Caddy, or Traefik on a VPS

## 🔒 Security Notes

### Best Practices

1. **HTTPS Only**: Always use HTTPS for both incoming (client → Netlify) and outgoing (Netlify → origin) traffic. This is enforced by default.

2. **Hide Origin IP**: Never expose your origin server IP (`ad.sdupdates.news`) publicly. Only share the Netlify domain.

3. **Firewall Rules**: Consider whitelisting only Netlify's IP ranges on your origin server for additional security.

4. **Sensitive Logs**: Disable verbose logging in production to avoid leaking sensitive information.

5. **Environment Variables**: For production deployments, consider storing the origin URL in a Netlify environment variable:
   ```javascript
   const ORIGIN_URL = process.env.ORIGIN_URL || 'https://ad.sdupdates.news';
   ```

6. **Rate Limiting**: Implement rate limiting on your origin server to prevent abuse.

### Security Advantages

- ✅ **IP Masking**: Your real server IP is hidden from clients
- ✅ **DDoS Protection**: Netlify's CDN provides basic DDoS mitigation
- ✅ **TLS Termination**: Netlify handles SSL/TLS certificates automatically
- ✅ **Geographic Distribution**: Requests are served from edge nodes globally

### Security Limitations

- ⚠️ Netlify can see your traffic (they are the middle layer)
- ⚠️ Not suitable for maximum-security scenarios requiring end-to-end encryption beyond TLS
- ⚠️ Origin server must still be secured properly

## 🛠️ Customization

### Change Origin Server

To use a different origin server, edit `functions/proxy.js`:

```javascript
const ORIGIN_URL = 'https://your-new-server.com';
```

Also update the Host header:
```javascript
forwardHeaders['Host'] = 'your-new-server.com';
```

### Use Environment Variables

For better security, use Netlify environment variables:

1. In Netlify Dashboard, go to **Site settings → Environment variables**

2. Add a variable:
   - **Key**: `ORIGIN_URL`
   - **Value**: `https://ad.sdupdates.news`

3. Update `functions/proxy.js`:
   ```javascript
   const ORIGIN_URL = process.env.ORIGIN_URL;
   ```

4. Redeploy your site

### Custom Domain

To use a custom domain instead of `*.netlify.app`:

1. In Netlify Dashboard, go to **Domain settings → Add custom domain**

2. Add your domain (e.g., `proxy.yourdomain.com`)

3. Configure DNS:
   - Add a CNAME record pointing to your Netlify site
   - Or use Netlify DNS for automatic configuration

4. Enable HTTPS (automatic with Let's Encrypt)

5. Update your V2Ray client configuration with the new domain

## 📊 Monitoring

### View Function Logs

```bash
netlify functions:log proxy
```

Or view logs in the Netlify Dashboard under **Functions → proxy → View logs**

### Common Issues

**502 Bad Gateway**
- Origin server is down or unreachable
- Firewall blocking Netlify IPs
- SSL/TLS certificate issues on origin

**504 Gateway Timeout**
- Request took longer than 10 seconds
- Origin server is slow to respond
- Consider upgrading to Pro plan for 26-second timeout

**403 Forbidden**
- Origin server blocking requests
- Check Host header and User-Agent
- Verify origin server allows the IP range

## 📝 Technical Implementation Details

### How the Proxy Function Works

1. **Request Parsing**: Extracts path, query parameters, and headers from incoming request

2. **Header Filtering**: Removes Netlify-specific headers and adds proper Host header for origin

3. **Request Forwarding**: Uses `fetch()` API to forward request to origin with all original parameters

4. **Response Handling**: Reads response as binary data, encodes as base64, and returns with original headers

5. **Error Handling**: Catches connection errors and returns proper HTTP error codes

### Key Features

- ✅ Preserves HTTP method (GET, POST, PUT, DELETE, etc.)
- ✅ Forwards all query parameters
- ✅ Preserves request/response headers
- ✅ Supports binary data (base64 encoding)
- ✅ Handles errors gracefully
- ✅ Logs requests for debugging

## 📚 Additional Resources

- [Netlify Functions Documentation](https://docs.netlify.com/functions/overview/)
- [V2Ray Documentation](https://www.v2ray.com/)
- [Netlify CLI Reference](https://cli.netlify.com/)
- [Netlify Environment Variables](https://docs.netlify.com/environment-variables/overview/)

## 🤝 Support

For issues specific to:
- **Netlify deployment**: Check Netlify documentation or support
- **V2Ray configuration**: Refer to V2Ray community resources
- **This proxy setup**: Review function logs and ensure origin server is accessible

## 📄 License

This project configuration is provided as-is for educational and personal use.

---

**⚡ Ready to Deploy!** Follow the deployment instructions above to get your Netlify reverse proxy running.
