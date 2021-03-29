package impl

import (
	"encoding/json"
	"fmt"

	"github.com/kachmazoff/doit-api/internal/model"
	"github.com/kachmazoff/doit-api/internal/repository"
)

type TimelineService struct {
	repo repository.Timeline

	srvChallenges   ChallengesService
	srvNotes        NotesService
	srvParticipants ParticipantsService
	srvSuggestions  SuggestionsService
}

func NewTimelineService(
	repo repository.Timeline,
	srvChallenges ChallengesService,
	srvNotes NotesService,
	srvParticipants ParticipantsService,
	srvSuggestions SuggestionsService,
) *TimelineService {
	return &TimelineService{
		repo:            repo,
		srvChallenges:   srvChallenges,
		srvNotes:        srvNotes,
		srvParticipants: srvParticipants,
		srvSuggestions:  srvSuggestions,
	}
}

func (s *TimelineService) GetAll() ([]model.TimelineItem, error) {
	return s.repo.GetAll()
}

func (s *TimelineService) GetCommon() ([]model.TimelineItem, error) {
	timeline, err := s.repo.GetCommon()
	return s.commonGetterHandler(timeline, err, true, true)
}

func (s *TimelineService) GetForUser(userId string) ([]model.TimelineItem, error) {
	timeline, err := s.repo.GetForUser(userId)
	return s.commonGetterHandler(timeline, err, true, true)
}

func (s *TimelineService) GetUserOwn(userId string) ([]model.TimelineItem, error) {
	timeline, err := s.repo.GetUserOwn(userId)
	return s.commonGetterHandler(timeline, err, true, false)
}

func (s *TimelineService) commonGetterHandler(timeline []model.TimelineItem, err error, needEnrich, needAnonymize bool) ([]model.TimelineItem, error) {
	if err != nil {
		return []model.TimelineItem{}, err
	}

	if needEnrich {
		// TODO: handle error
		s.EnrichTimeline(&timeline)
	}

	if needAnonymize {
		s.Anonymize(&timeline)
	}

	return timeline, nil
}

func (s *TimelineService) EnrichTimeline(timeline *[]model.TimelineItem) error {
	var err error

	for i := 0; i < len(*timeline); i++ {
		currErr := s.EnrichItem(&(*timeline)[i])
		if currErr != nil {
			err = currErr
		}
	}

	return err
}

func (s *TimelineService) EnrichItem(timelineItem *model.TimelineItem) error {
	// TODO: refactoring

	if timelineItem.ParticipantId != nil {
		participant, err := s.srvParticipants.GetById(*timelineItem.ParticipantId)

		if err != nil {
			return err
		}

		timelineItem.Participant = &participant
	}

	if timelineItem.NoteId != nil {
		note, err := s.srvNotes.GetById(*timelineItem.NoteId)

		if err != nil {
			return err
		}

		timelineItem.Note = &note
	}

	if timelineItem.SuggestionId != nil {
		suggestion, err := s.srvSuggestions.GetById(*timelineItem.SuggestionId)

		if err != nil {
			return err
		}

		timelineItem.Suggestion = &suggestion
	}

	return nil
}

func (s *TimelineService) Anonymize(timeline *[]model.TimelineItem) bool {
	isAnonym := false
	for i := 0; i < len(*timeline); i++ {
		isAnonym = isAnonym || s.AnonymizeItem(&(*timeline)[i])
	}
	return isAnonym
}

func (s *TimelineService) AnonymizeItem(timelineItem *model.TimelineItem) bool {
	isAnonym := false

	// TODO: refactoring
	isAnonym = isAnonym || s.srvChallenges.Anonymize(&timelineItem.Challenge)
	isAnonym = isAnonym || timelineItem.ParticipantId != nil && timelineItem.Participant != nil && timelineItem.Participant.Anonymous
	isAnonym = isAnonym || timelineItem.NoteId != nil && timelineItem.Note != nil && *timelineItem.Note.Anonymous
	isAnonym = isAnonym || timelineItem.SuggestionId != nil && timelineItem.Suggestion != nil && timelineItem.Suggestion.Anonymous

	if isAnonym {
		timelineItem.UserId = ""
		timelineItem.User = nil
	}

	return isAnonym
}

func print(timelineItem model.TimelineItem) {
	b, err := json.Marshal(timelineItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
