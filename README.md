# Bringy - a Telegram bot
<div align="center">
  <img src="Bringy_logo.jpg" width="200">
  <p><i>Логотип. Может быть, не самый красивый</i></p>
  <p><h3><a href="https://t.me/bringy_bot">@Bringy_bot</a></h3></p>
</div>
<hr>
Bringy позволяет быстро узнавать, о чём беседуют участники в Телеграм-группе. Отличительной особенностью является возможность отдельной работы по топикам, так что несколько тем не будут смешиваться в одну. Бот написан с помощью языка Go и библиотеки <a href="https://github.com/go-telegram/bot">go-telegram/bot</a>. Анализ сообщений участников происходит через модель Gemini, токен для которой предоставляется администраторами каждой из групп отдельно (собственный). В качестве базы данных используется MongoDB.
<hr>
В <code>services/config/config.go</code> можно найти файл, в котором можно настроить следующие параметры:
<ul>
  <li><code>MinuteIntervalSummarization</code> - временные промежутки в минутах между отправлением сообщений нейросети</li>
  <li><code>UTCPlusHours</code> - для редактирования UTC-зоны, которую использует бот. По умолчанию UTC+3 (Время МСК)</li>
  <li><code>ModelName</code> - название используемой модели Gemini</li>
  <li><code>ModelTemperature</code> - температура, с которой модель обрабатывает запрос</li>
</ul>
<hr>
Для запуска бота потребуется установленный компилятор Go, а также некоторые env параметры:
<ul>
  <li><code>BotToken</code> - токен бота Telegram</li>
  <li><code>MongoDBAppName</code> - оно указывается при создании базы данных в MongoDB</li>
  <li><code>MongoDBDatabaseName</code> - название базы данных MongoDB</li>
  <li><code>MongoDBUsername</code> - имя пользователя для БД</li>
  <li><code>MongoDBPassword</code> - пароль пользователя</li>
</ul>
