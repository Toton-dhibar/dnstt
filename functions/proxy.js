/**
 * Netlify Function: V2Ray xhttp + TLS Reverse Proxy
 * 
 * This function acts as a reverse proxy, forwarding all requests
 * to the origin V2Ray server while preserving headers, methods,
 * and body content.
 */

const ORIGIN_URL = 'https://ad.sdupdates.news';

exports.handler = async (event, context) => {
  try {
    // Extract the path after /xhttp from the Netlify function invocation
    const path = event.path.replace('/.netlify/functions/proxy', '');
    
    // Construct the full target URL
    const targetUrl = `${ORIGIN_URL}/xhttp${path}${event.rawQuery ? '?' + event.rawQuery : ''}`;
    
    console.log(`[Proxy] ${event.httpMethod} ${targetUrl}`);
    
    // Prepare headers to forward (exclude Netlify-specific headers)
    const forwardHeaders = {};
    const excludeHeaders = [
      'host',
      'x-forwarded-for',
      'x-forwarded-proto',
      'x-forwarded-host',
      'x-netlify-id',
      'client-ip',
      'true-client-ip'
    ];
    
    for (const [key, value] of Object.entries(event.headers)) {
      const lowerKey = key.toLowerCase();
      if (!excludeHeaders.includes(lowerKey)) {
        forwardHeaders[key] = value;
      }
    }
    
    // Set the correct Host header for the origin
    forwardHeaders['Host'] = 'ad.sdupdates.news';
    
    // Prepare fetch options
    const fetchOptions = {
      method: event.httpMethod,
      headers: forwardHeaders,
    };
    
    // Add body for methods that support it
    if (event.httpMethod !== 'GET' && event.httpMethod !== 'HEAD' && event.body) {
      fetchOptions.body = event.isBase64Encoded 
        ? Buffer.from(event.body, 'base64')
        : event.body;
    }
    
    // Forward the request to the origin
    const response = await fetch(targetUrl, fetchOptions);
    
    // Extract response headers
    const responseHeaders = {};
    response.headers.forEach((value, key) => {
      // Skip headers that might cause issues
      const lowerKey = key.toLowerCase();
      if (lowerKey !== 'transfer-encoding' && lowerKey !== 'connection') {
        responseHeaders[key] = value;
      }
    });
    
    // Read the response body
    const responseBody = await response.arrayBuffer();
    const body = Buffer.from(responseBody).toString('base64');
    
    // Return the proxied response
    return {
      statusCode: response.status,
      headers: responseHeaders,
      body: body,
      isBase64Encoded: true
    };
    
  } catch (error) {
    console.error('[Proxy Error]', error.message);
    
    // Return error response
    return {
      statusCode: 502,
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        error: 'Bad Gateway',
        message: 'Failed to connect to origin server',
        details: error.message
      })
    };
  }
};
