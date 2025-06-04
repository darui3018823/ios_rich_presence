# 2025 iOS ShortCut DiscordRPC Server: darui3018823 All rights reserved.
# All works created by darui3018823 associated with this repository are the intellectual property of darui3018823.
# Packages and other third-party materials used in this repository are subject to their respective licenses and copyrights.

# clear_rpc.py
from pypresence import Presence
import json
import sys
import os

def load_rpc_config(app_name: str) -> dict:
    json_path = f"./json/{app_name}.json"
    if not os.path.exists(json_path):
        print(f"設定ファイルが見つかりません: {json_path}")
        sys.exit(1)
    with open(json_path, encoding="utf-8") as f:
        return json.load(f).get(app_name, {})

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: clear_rpc.py <AppName>")
        sys.exit(1)

    app_name = sys.argv[1]
    config = load_rpc_config(app_name)
    client_id = config.get("ClientId")

    if not client_id:
        print("ClientId が設定ファイルに存在しません")
        sys.exit(1)

    RPC = Presence(client_id)
    RPC.connect()
    RPC.clear()  # RPC解除
    print(f"RPC cleared for {app_name}")
