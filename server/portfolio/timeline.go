package portfolio

import db "github.com/matejgrzinic/portfolio/db/queries"

func (p *Portfolio) GetUserTimeline(user string, timeframe string) (*[]db.DbTimelineData, error) {
	timeline, err := p.db.Query.GetUserTimeline(user, timeframe)
	if err != nil {
		return nil, err
	}
	return timeline, nil
}
