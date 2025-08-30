import os
import subprocess
import signal
import atexit
import requests
from flask import Flask, request, Response, render_template, jsonify, redirect, stream_with_context

# Create a Flask app that proxies to the Go server
app = Flask(__name__, static_folder=None)  # Disable Flask's static folder handling
app.secret_key = os.environ.get("SESSION_SECRET", "your-secret-key")

# Go server configuration
GO_PORT = "8080"
GO_SERVER_URL = f"http://localhost:{GO_PORT}"

# Start the Go server
go_process = None

def start_go_server():
    global go_process
    try:
        print("Starting Go server...")
        # Start the Go server as a child process with the specified port
        go_process = subprocess.Popen(["go", "run", "main.go"], 
                                     stdout=subprocess.PIPE, 
                                     stderr=subprocess.PIPE,
                                     env=dict(os.environ, GO_PORT=GO_PORT))
        print(f"Go server started with PID: {go_process.pid}")
    except Exception as e:
        print(f"Error starting Go server: {e}")
        go_process = None

def cleanup_go_server():
    global go_process
    if go_process:
        print(f"Terminating Go server with PID: {go_process.pid}")
        try:
            go_process.terminate()
            go_process.wait(timeout=5)
        except:
            go_process.kill()

# Register the cleanup function to run when the app exits
atexit.register(cleanup_go_server)

# Start the Go server when the Flask app starts
start_go_server()

# Proxy all routes to the Go server
@app.route('/', defaults={'path': ''}, methods=['GET', 'POST', 'PUT', 'DELETE', 'PATCH'])
@app.route('/<path:path>', methods=['GET', 'POST', 'PUT', 'DELETE', 'PATCH'])
def proxy(path):
    try:
        # Forward the request to the Go server
        url = f"{GO_SERVER_URL}/{path}"
        
        # Copy request headers
        headers = {key: value for key, value in request.headers if key != 'Host'}
        
        # Special handling for JSON content
        json_data = None
        if request.is_json:
            json_data = request.get_json()
            
        # Make the request to the Go server
        resp = requests.request(
            method=request.method,
            url=url,
            headers=headers,
            params=request.args,
            json=json_data if request.is_json else None,
            data=None if request.is_json else request.get_data(),
            cookies=request.cookies,
            allow_redirects=False,
            stream=True
        )
        
        # Create Flask response
        response = Response(
            stream_with_context(resp.iter_content(chunk_size=1024)),
            status=resp.status_code
        )
        
        # Copy response headers
        for key, value in resp.headers.items():
            if key.lower() not in ('content-length', 'transfer-encoding', 'content-encoding'):
                response.headers[key] = value
                
        return response
    except requests.RequestException as e:
        return f"Error connecting to Go server: {e}", 500
