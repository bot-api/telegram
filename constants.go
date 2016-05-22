package telegram

// Type of chat
const (
	PrivateChatType    = "private"
	GroupChatType      = "group"
	SuperGroupChatType = "supergroup"
	ChannelChatType    = "channel"
)

// Type of action to broadcast.
//
// Choose one, depending on what the user is about to receive:
// 	typing for text messages
// 	upload_photo for photos
// 	record_video or upload_video for videos
// 	record_audio or upload_audio for audio files
//	upload_document for general files
// 	find_location for location data
const (
	ActionTyping         = "typing"
	ActionUploadPhoto    = "upload_photo"
	ActionRecordVideo    = "record_video"
	ActionUploadVideo    = "upload_video"
	ActionRecordAudio    = "record_audio"
	ActionUploadAudio    = "upload_audio"
	ActionUploadDocument = "upload_document"
	ActionFindLocation   = "find_location"
)

// internal constants for method names
const (
	getMeMethod                 = "getMe"
	getUpdatesMethod            = "getUpdates"
	getUserProfilePhotosMethod  = "getUserProfilePhotos"
	getChatMethod               = "getChat"
	getChatAdministratorsMethod = "getChatAdministrators"
	getChatMembersCountMethod   = "getChatMembersCount"
	getChatMemberMethod         = "getChatMember"
	kickChatMemberMethod        = "kickChatMember"
	unbanChatMemberMethod       = "unbanChatMember"
	leaveChatMethod             = "leaveChat"

	sendChatActionMethod = "sendChatAction"
	sendMessageMethod    = "sendMessage"
	sendVenueMethod      = "sendVenue"
	sendPhotoMethod      = "sendPhoto"
	sendAudioMethod      = "sendAudio"
	sendVideoMethod      = "sendVideo"
	sendVoiceMethod      = "sendVoice"
	sendDocumentMethod   = "sendDocument"
	sendContactMethod    = "sendContact"
	sendLocationMethod   = "sendLocation"
	sendStickerMethod    = "sendSticker"
	forwardMessageMethod = "forwardMessage"

	answerCallbackQueryMethod = "answerCallbackQuery"
	setWebhookMethod          = "setWebhook"
	getFileMethod             = "getFile"
	answerInlineQueryMethod   = "answerInlineQuery"

	editMessageTextMethod        = "editMessageText"
	editMessageCaptionMethod     = "editMessageCaption"
	editMessageReplyMarkupMethod = "editMessageReplyMarkup"
)

// constants for field names for file-like messages
const (
	photoField    = "photo"
	documentField = "document"
	audioField    = "audio"
	stickerField  = "sticker"
	videoField    = "video"
	voiceField    = "voice"
)

// Constant values for ParseMode in MessageCfg.
const (
	MarkdownMode = "Markdown"
	HTMLMode     = "HTML"
)

// EntityType constants helps to set type of entity in MessageEntity object
const (
	// @username
	MentionEntityType     = "mention"
	HashTagEntityType     = "hashtag"
	BotCommandEntityType  = "bot_command"
	URLEntityType         = "url"
	EmailEntityType       = "email"
	BoldEntityType        = "bold"         // bold text
	ItalicEntityType      = "italic"       // italic text
	CodeEntityType        = "code"         // monowidth string
	PreEntityType         = "pre"          // monowidth block
	TextLinkEntityType    = "text_link"    // for clickable text URLs
	TextMentionEntityType = "text_mention" // for users without usernames
)

// ChatMember possible statuses
const (
	MemberStatus              = "member"
	CreatorMemberStatus       = "creator"
	AdministratorMemberStatus = "administrator"
	LeftMemberStatus          = "left"
	KickedMemberStatus        = "kicked"
)
