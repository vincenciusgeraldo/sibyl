package common

import (
	"github.com/vincenciusgeraldo/sibyl/pkg/models"
	"fmt"
	"github.com/vincenciusgeraldo/sibyl/pkg/utils"
)

func BuildReview(review models.Review, status models.PullRequestStatus) string {
	message := ""
	prName := review.PRName
	prFormat := "`%s#%d` - `%s`"
	if prName == "" {
		prFormat = "`%s#%d`"
	}

	if review.Emergency {
		message += "\xE2\x80\xBC *Emergency* \xE2\x80\xBC\n" +
			fmt.Sprintf("Antrian PR "+ prFormat +" dari %s memerlukan _review_ secepatnya, PR-nya bisa dilihat [disini](https://github.com/bukalapak/%s/pull/%d) ya.\n", review.Repo, review.PRNumber, prName, review.Requester, review.Repo, review.PRNumber)
	} else {
		message += fmt.Sprintf("Antrian _review_ PR "+ prFormat +" dari %s. PR-nya bisa dilihat [disini](https://github.com/bukalapak/%s/pull/%d) ya.\n", review.Repo, review.PRNumber, prName, review.Requester, review.Repo, review.PRNumber)
	}

	message += "`Check Status : `\n"
	if !status.Mergeable && status.MergeableStatus == "dirty" {
		message += "- This PR has unresolved conflict!\n"
	} else {
		for _, c := range status.Statuses {
			message += fmt.Sprintf("- [%s](%s) - %s\n", c.GetContext(), c.GetURL(), c.GetDescription())
		}

		if len(status.Statuses) == 0 {
			message += "- No available check statuses for this PR!\n"
		}
	}

	message += "`Review Status : `\n"
	for _, r := range review.Reviewers {
		rst := ""
		if utils.ArrayInclude(review.ApprovedBy, r) {
			rst = "approved \xF0\x9F\x91\x8D"
		} else if utils.ArrayInclude(review.ReviewedBy, r) {
			rst = "need changes"
		} else {
			rst = "not reviewed"
		}
		message += fmt.Sprintf("- %s - %s\n", r, rst)
	}

	return message + "\n"
}
