package wenda

import "github.com/kongchuanhujiao/server/internal/app/datahub/public/wenda"

// CalculateQuestion 计算问答结果
func CalculateQuestion(w *wenda.Detail) (calc *wenda.CalculationsTab) {

	var (
		rightStus []uint64
		wrongStus []wenda.CalculationsWrong
	)

	rightStus = []uint64{}
	wrongStus = []wenda.CalculationsWrong{}

	correctAnswer := w.Questions.Key

	for _, ans := range w.Answers {
		if ans.Answer == correctAnswer {
			rightStus = append(rightStus, ans.QQ)
		} else {
			status := getWrongStatusByOption(wrongStus, ans.Answer)

			if status == nil {
				wrongStus = append(wrongStus, wenda.CalculationsWrong{
					Type:   ans.Answer,
					Member: []uint64{ans.QQ},
				})
			} else {
				status.Member = append(status.Member, ans.QQ)
			}
		}
	}

	calc = &wenda.CalculationsTab{
		Question: w.Questions.ID,
		Count:    uint8(len(w.Answers)),
		Right:    rightStus,
		Wrong:    wrongStus,
	}

	return
}

func getWrongStatusByOption(wrongStus []wenda.CalculationsWrong, option string) *wenda.CalculationsWrong {
	for _, status := range wrongStus {
		if status.Type == option {
			return &status
		}
	}

	return nil
}
