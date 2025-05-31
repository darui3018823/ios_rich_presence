# 2025 iOS ShortCut DiscordRP: darui3018823 All rights reserved.
# All works created by darui3018823 associated with this repository are the intellectual property of darui3018823.
# Packages and other third-party materials used in this repository are subject to their respective licenses and copyrights.

import sys
import json
import time
from pypresence import Presence

# 引数取得
app = sys.argv[1]
device = sys.argv[2] if len(sys.argv) > 2 else "Unknown Device"
user = sys.argv[3] if len(sys.argv) > 3 else "Unknown User"

print(f"DEBUG: 読み込みファイル: ./json/{app}.json")
print("DEBUG: 引数", app, device, user)

# JSON読み込み
with open(f"./json/{app}.json", "r", encoding="utf-8") as f:
    all_config = json.load(f)

cfg = all_config.get(app)
if not cfg:
    print(f"Error: app config for '{app}' not found.")
    sys.exit(1)
    
print("DEBUG: 読み込んだボタン:", cfg.get("Buttons"))

# プレースホルダ置換
def replace_placeholders(text: str, context: dict):
    for key, value in context.items():
        text = text.replace(f"{{{key}}}", value)
    return text

context = {
    "device": device,
    "user": user
}

for key in ["Details", "State", "LargeImageText", "SmallImageText"]:
    if key in cfg and isinstance(cfg[key], str):
        cfg[key] = replace_placeholders(cfg[key], context)

# Discord RPC 初期化と更新
client_id = cfg["ClientId"]
RPC = Presence(client_id)
RPC.connect()

kwargs = {
    "state": cfg.get("State"),
    "details": cfg.get("Details"),
    "large_image": cfg.get("LargeImage"),
    "large_text": cfg.get("LargeImageText"),
    "small_image": cfg.get("SmallImage"),
    "small_text": cfg.get("SmallImageText"),
    "buttons": cfg.get("Buttons"),
    "party_id": cfg.get("PartyId"),
    "party_size": cfg.get("PartySize"),
    "start": time.time()
}

# None を除去して実行
RPC.update(**{k: v for k, v in kwargs.items() if v is not None})
print(f"[OK] RPC updated for app: {app} ({device}, {user})")

print("RPC set. Holding forever...")
while True:
    time.sleep(3600)