package message

import (
	"fmt"
	"wpp-cloud/internal/tracking"
)

func GetGenericTrackError(code string) string {
	return fmt.Sprintf("Pedimos desculpa mas no momento não foi possível "+
		"completar o rastreio do código %s. Tente novamente dentro de alguns minutos.\n"+
		"Caso o problema persista, entre em contato com o administrador.", code)
}

func GetHelp() string {
	return fmt.Sprintf("%s%s%s", GetGreetings(), GetCommandsAvailable(), GetContact())
}

func GetContact() string {
	return "\nEsse BOT é desenvolvido e mantido por Mateus Matinato.\n" +
		"Linkedin: https://www.linkedin.com/in/mateusmatinato/\n" +
		"Github: https://github.com/mateusmatinato\n" +
		"Contato profissional: mateusmatinato@gmail.com\n"
}

func GetCommandsAvailable() string {
	return "\nNo momento temos os seguintes comandos disponíveis:\n" +
		"*/rastrear*: Rastreie seus pacotes;\n" +
		"*/ajuda*: Veja os comandos disponíveis e mais informações;\n\n" +
		"Para entender o funcionamento de cada comando, basta enviar uma mensagem com o comando desejado.\n"
}

func GetTrackUpdateSuccess(code string, result tracking.TrackPackageResult) string {
	updateMsg := "Tem atualização no seu pacote!!!"
	return fmt.Sprintf("%s\n%s", updateMsg, GetTrackSuccess(code, result))
}

func GetTrackSuccess(code string, result tracking.TrackPackageResult) string {
	return fmt.Sprintf("Segue informações do seu pacote _%s_:\n"+
		"*Status*: %s\n"+
		"*Local*: %s\n"+
		"*Horário de Atualização*: %s\n",
		code, result.Status, result.Place, fmt.Sprintf("%s - %s", result.Date, result.Time))
}

func GetGenericMessageError() string {
	return "Não consegui entender sua mensagem. Estou em fase de testes e caso você tenha dúvidas " +
		"de como usar o *TrackBOT* basta enviar o comando _/ajuda_."
}

func GetGreetings() string {
	return "Bem vindo ao *TrackBOT*!\n"
}

func GetTrackSteps() string {
	return "Para rastrear suas encomendas basta enviar uma " +
		"mensagem no seguinte formato: \n" +
		"_/rastrear CODIGO_RASTREIO_\n\n" +
		"Por exemplo: \n" +
		"_/rastrear NA209663501BR_\n\n" +
		"Observação: Só é possível rastrear *UM PACOTE* por vez."
}
