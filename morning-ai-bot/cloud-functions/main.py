import datetime
import json
import os

import requests
from linebot import LineBotApi
from linebot.models import TextSendMessage
from pytz import timezone

#===== common =====
class MyException(Exception):
    pass

# https://developer.twitter.com/en/docs/tweets/timelines/api-reference/get-statuses-user_timeline
def search_oldest_tweet(token, name, from_jst=None):
    """
    from_jst... yyyymmddhhmm
    """
    url = "https://api.twitter.com/1.1/statuses/user_timeline.json"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    payload = {
        "count": 10,
        "screen_name": name,
        "trim_user": "t"
    }
    if from_jst is not None:
        from_datetime = datetime.datetime.strptime(from_jst+"+0900", "%Y%m%d%H%M%z")
    else:
        from_tmp = datetime.datetime.now(timezone("Asia/Tokyo"))
        from_datetime = timezone("Asia/Tokyo").localize(datetime.datetime(
            year=from_tmp.year,
            month=from_tmp.month,
            day=from_tmp.day,
            hour=5,
            minute=55
        ))
    res = requests.get(url, headers=headers, params=payload)
    res_json = json.loads(res.text)
    res_limitted = [
        x
        for x in res_json
        if from_datetime <= datetime.datetime.strptime(x["created_at"], "%a %b %d %H:%M:%S %z %Y")
    ]
    if len(res_limitted) == 0:
        raise MyException(f"couldn't find any tweets\n{res.text}")
    else:
        tweet_id = res_limitted[-1]["id_str"]
        # results is automatically sorted in reverse chronological order
        tweet_url = f"https://twitter.com/aichan_nel/status/{tweet_id}"
        return tweet_url

def send_message(message):
    line_bot_api = LineBotApi(os.getenv("MORNING_AI_BOT_TOKEN"))
    line_bot_api.broadcast(TextSendMessage(text=message))

def morning_ai_bot(_):
    tweet_url = search_oldest_tweet(
        os.getenv("TWITTER_TOKEN"),
        'aichan_nel',
    )
    send_message(tweet_url)
