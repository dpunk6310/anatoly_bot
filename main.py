import sys
import logging
import asyncio

from aiogram import Bot, Dispatcher, types
from aiogram import Router


API_TOKEN = ''

bot = Bot(token=API_TOKEN)
dp = Dispatcher()

SOURCE_CHANNEL_ID = -100
TARGET_CHANNEL_ID = -100

router = Router()

@router.channel_post()
async def handle_new_post(message: types.Message):
    print(message.chat.id)
    if message.chat.id == SOURCE_CHANNEL_ID:
        if message.text:
            await bot.send_message(chat_id=TARGET_CHANNEL_ID, text=message.text)
        if message.photo:
            await bot.send_photo(chat_id=TARGET_CHANNEL_ID, photo=message.photo[-1].file_id, caption=message.caption)
        if message.video:
            await bot.send_video(chat_id=TARGET_CHANNEL_ID, video=message.video.file_id, caption=message.caption)
        if message.document:
            await bot.send_document(chat_id=TARGET_CHANNEL_ID, document=message.document.file_id, caption=message.caption)

dp.include_router(router)

if __name__ == '__main__':
    
    logging.basicConfig(level=logging.INFO, stream=sys.stdout)
    asyncio.run(dp.start_polling(bot), debug=True)