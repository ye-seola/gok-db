import asyncio
from websockets.asyncio.client import ClientConnection, connect as ws_connect
import json

RECONNECT_INTERVAL = 1


async def on_message(ws: ClientConnection, message):
    try:
        data = json.loads(message)

        if "error" in data:
            print(data)
            return

        event = data.get("event")
        payload = data.get("payload", {})

        async def reply(msg: str, chat_id):
            await ws.send(
                json.dumps(
                    {
                        "action": "SENDMSG",
                        "payload": {"message": str(msg), "chatId": chat_id},
                    }
                )
            )

        if event == "MSG":
            msg = payload.get("message")
            print(data)

            if msg == "/Hello":
                await reply("Hello")

    except Exception as e:
        print(f"Error handling message: {e}")


async def connect_ws():
    ws_url = "ws://127.0.0.1:9023/ws"
    while True:
        try:
            async with ws_connect(ws_url) as websocket:
                print("Connected to WebSocket server")
                async for message in websocket:
                    try:
                        await on_message(websocket, message)
                    except Exception:
                        print("on_message error")
        except Exception as e:
            print(f"An error occurred: {e}.")
        await asyncio.sleep(RECONNECT_INTERVAL)


if __name__ == "__main__":
    asyncio.run(connect_ws())
