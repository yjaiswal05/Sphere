# Email Domain Verifier

A web-based tool that validates email domain configurations by checking DNS records including MX, SPF, and DMARC records. The application provides a user-friendly interface for domain verification tasks and displays comprehensive results about email security configurations.

![Email Domain Verifier](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![Python](https://img.shields.io/badge/Python-3776AB?style=flat&logo=python&logoColor=white)
![Flask](https://img.shields.io/badge/Flask-000000?style=flat&logo=flask&logoColor=white)
![Bootstrap](https://img.shields.io/badge/Bootstrap-7952B3?style=flat&logo=bootstrap&logoColor=white)

## Features

- **MX Record Validation**: Checks if mail exchange servers are properly configured
- **SPF Record Analysis**: Validates Sender Policy Framework settings
- **DMARC Policy Verification**: Examines Domain-based Message Authentication policies
- **Responsive Web Interface**: Clean, mobile-friendly design with dark theme
- **Real-time Results**: Instant verification with detailed technical information
- **Error Handling**: Graceful handling of DNS lookup failures and network issues

## Architecture

This project uses a hybrid architecture:

- **Go Backend**: Core verification engine that performs DNS lookups and domain analysis
- **Flask Proxy**: Lightweight Python proxy layer for compatibility with various deployment platforms
- **Static Frontend**: HTML5, CSS3, and JavaScript with Bootstrap for responsive design

The Go server runs on port 8080 and handles all verification logic, while the Flask application on port 5000 acts as a proxy to forward requests and manage the Go server lifecycle.

## Prerequisites

Before running this project, ensure you have:

- **Go** (version 1.19 or later) - [Download here](https://golang.org/dl/)
- **Python** (version 3.8 or later) - [Download here](https://python.org/downloads/)
- **Git** (for cloning the repository) - [Download here](https://git-scm.com/)

## Local Installation & Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd email-domain-verifier
```

### 2. Install Python Dependencies

```bash
# Using pip
pip install flask gunicorn requests

# Or using pip with requirements
pip install -r heroku_requirements.txt
```

### 3. Verify Go Installation

```bash
# Check if Go is installed
go version

# If not installed, download from https://golang.org/dl/
```

## Running the Application

###  Using Flask Development Server 

```bash
# Start the application
python -m flask run --host=0.0.0.0 --port=5000

# The Flask app will automatically start the Go server
# Access the application at: http://localhost:5000
```

## How to Use

1. **Access the Web Interface**: Open your browser and navigate to `http://localhost:5000`
2. **Enter Domain**: Type a domain name (e.g., "google.com", "github.com") in the input field
3. **Verify**: Click the "Verify Domain" button
4. **View Results**: The application will display:
   - MX Records status (mail server configuration)
   - SPF Record status and content (sender authorization)
   - DMARC Record status and content (email authentication policy)

## Example Domains to Test

- `google.com` - Well-configured domain with all records
- `github.com` - Popular service with email security
- `example.com` - Basic configuration for testing
- `your-domain.com` - Test your own domain

## API Endpoints

### POST /verify
- **Description**: Verifies domain email configuration
- **Content-Type**: application/json
- **Request Body**:
  ```json
  {
    "domain": "example.com"
  }
  ```
- **Response**:
  ```json
  {
    "domain": "example.com",
    "hasMX": true,
    "hasSPF": true,
    "hasDMARC": true,
    "spfRecord": "v=spf1 -all",
    "dmarcRecord": "v=DMARC1;p=reject;sp=reject;adkim=s;aspf=s"
  }
  ```
## Troubleshooting

### Common Issues

1. **Go server not starting**
   - Ensure Go is installed: `go version`
   - Check port availability: `netstat -an | grep 8080`

2. **Python dependencies missing**
   - Install Flask: `pip install flask`
   - Install requests: `pip install requests`

3. **DNS lookup failures**
   - Check internet connection
   - Try different domains
   - Some corporate networks may block DNS queries

4. **Port conflicts**
   - Change the Go port: `export GO_PORT=8081`
   - Change the Flask port: `flask run --port=5001`

### Logs and Debugging

```bash
# View application logs (when running locally)
# Check terminal output for error messages

# For Heroku deployment
heroku logs --tail
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Technical Details

### DNS Verification Process

1. **MX Records**: Uses `net.LookupMX()` to find mail exchange servers
2. **SPF Records**: Searches TXT records for entries starting with `v=spf1`
3. **DMARC Records**: Queries `_dmarc.domain` for TXT records with `v=DMARC1`

### Security Features

- Input validation and sanitization
- Error handling for malformed requests
- Timeout protection for DNS queries
- No storage of user data or domains

---

**Built with ❤️ using Go, Python, and modern web technologies**