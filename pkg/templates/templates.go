package templates

const InternalServiceError = "" +
	"Упс, произошла ошибка внутри сервиса\n" +
	"\n" +
	"На данный момент я не смогу обработать твое обращение\n" +
	"\n" +
	"Повтори попытку через 5 минут"

const UsrIsNotKnownByTelegram = "" +
	"Извини, я еще не могу работать с тобой\n" +
	"Если ты дашь мне номер телефона - мы в расчете 😏"

const UsrIsKnownByTelegram = "" +
	"Я уже знаю тебя!\n" +
	"Ты можешь брать заказы в чате 🤟" +
	"\n" +
	"Попроси ссылку у Александра, если ты еще не там"

const UsrIsNotRunner = "" +
	"Спасибо за твои данные\n" +
	"Пока что ты не зарегистрирован как доставщик\n" +
	"Я вынужден отклонить твой запрос начать работу ❌"

const BeginWorkSuccess = "" +
	"Спасибо за твои данные\n" +
	"Я нашел тебя, %s. Теперь мы можем работать вместе!\n" +
	"Ты можешь брать заказы в чате 🤟\n" +
	"\n" +
	"Попроси ссылку у Александра, если ты еще не там"

const RunnerDoesNotExist = "" +
	"Курьер с номером %s\n" +
	"еще не зарегистрирован\n" +
	"\n" +
	"Прости 🙇"

const DeliveryHasAlreadyReserved = "" +
	"К сожалению, доставка %d уже зарезервирована!\n" +
	"Попробуй взять другую"

const Success = "✅"

const YouAreNotARunner = "" +
	"Вы не курьер! \U0001F976\n" +
	"Или вы пока еще не вошли как курьер...🙄\n" +
	"Хотите войти?\n"

const ComeHere = "Тебе сюда 😇"

const Report = "Сообщить о проблеме 🔰"

const Complete = "Доставил 🍀"

const AfterReserveReply = "\n" +
	"------------------------\n" +
	"Вы взяли этот заказ в %s\n" +
	"\n" +
	"Номер доставки #%d\n" +
	"\n" +
	"После доставки нажмите кнопку 🍀\n" +
	"*это не обязательно, просто чтобы вам было удобнее " +
	"следить за своими доставками!"
