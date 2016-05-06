# telegram

Implementation of the telegram bot API, inspired by github.com/go-telegram-bot-api/telegram-bot-api

Alpha version, don't use it.


# TODO:

- [ ] Add integration tests
- [ ] Add travis-ci integration
- [x] Handlers 
- [x] Command handlers
- [ ] Session handlers
- [x] Middleware
- [ ] Examples


- [ ] gopkg version
- [ ] documentation
- [ ] Benchmark ffjson and easyjson.
- [ ] GAE example. 
- [ ] Handle 
        status code: 409
        received: {"ok":false,"error_code":409,"description":"[Error]: Conflict: another webhook is active"}





# Supported API methods:
- [x] getMe
- [x] sendMessage
- [x] forwardMessage
- [x] sendPhoto
- [x] sendAudio
- [x] sendDocument
- [x] sendSticker
- [x] sendVideo
- [x] sendVoice
- [x] sendLocation
- [x] sendChatAction
- [x] getUserProfilePhotos
- [x] getUpdates
- [x] setWebhook
- [x] getFile
- [ ] answerInlineQuery inline bots

#  Supported API v2 methods:
- [x] sendVenue
- [x] sendContact
- [x] editMessageText
- [x] editMessageCaption
- [x] editMessageReplyMarkup
- [ ] kickChatMember
- [ ] unbanChatMember
- [x] answerCallbackQuery

# Supported Inline modes



Other bots:
I like this handler system
https://bitbucket.org/master_groosha/telegram-proxy-bot/src/07a6b57372603acae7bdb78f771be132d063b899/proxy_bot.py?fileviewer=file-view-default

