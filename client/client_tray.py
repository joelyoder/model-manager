import sys
import json
import os
import threading
import time
import requests
import websocket
from PIL import Image, ImageDraw
import pystray
from urllib.parse import urlparse

# Load Config
def get_application_path():
    if getattr(sys, 'frozen', False):
        return os.path.dirname(sys.executable)
    return os.path.dirname(os.path.abspath(__file__))

CONFIG_FILE = os.path.join(get_application_path(), 'config.json')
CONFIG = {}

def load_config():
    global CONFIG
    if os.path.exists(CONFIG_FILE):
        with open(CONFIG_FILE, 'r') as f:
            CONFIG = json.load(f)
    else:
        print(f"Config file not found at {CONFIG_FILE}")
        sys.exit(1)

load_config()

SERVER_URL = CONFIG.get('server_url')
API_KEY = CONFIG.get('api_key')
ROOT_PATH = CONFIG.get('root_path')
CLIENT_ID = CONFIG.get('client_id')

if not all([SERVER_URL, API_KEY, ROOT_PATH, CLIENT_ID]):
    print("Missing configuration values in config.json")
    sys.exit(1)

# Ensure root path exists
if not os.path.exists(ROOT_PATH):
    os.makedirs(ROOT_PATH)

ws = None
tray_icon = None

def create_icon():
    # Create a simple icon
    image = Image.new('RGB', (64, 64), (0, 128, 255))
    dc = ImageDraw.Draw(image)
    dc.rectangle((16, 16, 48, 48), fill=(255, 255, 255))
    return image

def on_message(ws_app, message):
    print(f"Received: {message}")
    try:
        data = json.loads(message)
        action = data.get('action')
        
        if action == 'download':
            handle_download(ws_app, data)
        elif action == 'delete':
            handle_delete(ws_app, data)
    except Exception as e:
        print(f"Error processing message: {e}")

def sanitize_path(subdirectory, filename):
    # Security check: Ensure path is within ROOT_PATH
    # Strip dangerous characters
    safe_filename = os.path.basename(filename)
    safe_subdir = subdirectory.replace('..', '').strip('/\\')
    
    full_path = os.path.join(ROOT_PATH, safe_subdir, safe_filename)
    full_path = os.path.abspath(full_path)
    
    if not full_path.startswith(os.path.abspath(ROOT_PATH)):
        raise ValueError("Path traversal attempt detected")
        
    return full_path

def handle_download(ws_app, data):
    url = data.get('url')
    model_version_id = data.get('model_version_id')

    # Resolve URL against server
    if url and url.startswith('/'):
        u = urlparse(SERVER_URL)
        scheme = 'https' if u.scheme == 'wss' else 'http'
        base_url = f"{scheme}://{u.netloc}"
        url = base_url + url

    try:
        # Determine target path from URL (remove /downloads/ prefix)
        rel_path = data.get('url', '').replace('/downloads/', '', 1)
        rel_path = rel_path.lstrip('/').lstrip('\\')
        
        from urllib.parse import unquote
        rel_path = unquote(rel_path)
        
        # Security check - basic
        if '..' in rel_path:
             raise ValueError("Invalid path components")
             
        target_path = os.path.join(ROOT_PATH, rel_path)
        target_path = os.path.abspath(target_path)
        
        if not target_path.startswith(os.path.abspath(ROOT_PATH)):
             raise ValueError("Path traversal attempt")
             
        os.makedirs(os.path.dirname(target_path), exist_ok=True)
        print(f"Downloading {url} to {target_path}...")
        
        # Download
        with requests.get(url, stream=True) as r:
            r.raise_for_status()
            with open(target_path, 'wb') as f:
                for chunk in r.iter_content(chunk_size=8192): 
                    f.write(chunk)
                    
        print("Download complete.")
        
        # Send confirmation
        response = {
            "type": "complete",
            "model_version_id": model_version_id
        }
        ws_app.send(json.dumps(response))
        
    except Exception as e:
        print(f"Download invalid: {e}")

def handle_delete(ws_app, data):
    filename = data.get('filename') # Now contains relative path
    model_version_id = data.get('model_version_id')
    
    try:
        # Security check
        if '..' in filename:
             raise ValueError("Invalid path components")

        target_path = os.path.join(ROOT_PATH, filename)
        target_path = os.path.abspath(target_path)
        
        if not target_path.startswith(os.path.abspath(ROOT_PATH)):
            raise ValueError("Path traversal attempt detected")

        if os.path.exists(target_path):
            os.remove(target_path)
            print(f"Deleted {target_path}")
            
            response = {
                "type": "deleted",
                "model_version_id": model_version_id
            }
            ws_app.send(json.dumps(response))
        else:
            print(f"File not found: {target_path}")
            
    except Exception as e:
        print(f"Delete failed: {e}")


def on_error(ws_app, error):
    print(f"WebSocket Error: {error}")

def on_close(ws_app, close_status_code, close_msg):
    print("WebSocket Closed")

def on_open(ws_app):
    print("WebSocket Connected")

def run_websocket():
    global ws
    headers = {
        "Authorization": API_KEY,
        "X-Client-ID": CLIENT_ID
    }
    
    while True:
        try:
            ws = websocket.WebSocketApp(SERVER_URL,
                                      header=headers,
                                      on_open=on_open,
                                      on_message=on_message,
                                      on_error=on_error,
                                      on_close=on_close)
            ws.run_forever()
        except Exception as e:
            print(f"Connection failed: {e}")
        
        print("Reconnecting in 5 seconds...")
        time.sleep(5)

def on_exit(icon, item):
    icon.stop()
    sys.exit(0)

def main():
    # Start WebSocket in a daemon thread
    t = threading.Thread(target=run_websocket)
    t.daemon = True
    t.start()
    
    # Run Tray Icon in main thread
    global tray_icon
    tray_icon = pystray.Icon("ModelClient", create_icon(), "Model Manager Client", menu=pystray.Menu(
        pystray.MenuItem("Client ID: " + CLIENT_ID, lambda: None, enabled=False),
        pystray.MenuItem("Exit", on_exit)
    ))
    tray_icon.run()

if __name__ == "__main__":
    main()
