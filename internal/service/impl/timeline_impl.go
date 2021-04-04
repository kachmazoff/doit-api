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

func (s *TimelineService) GetWithFilters(filters model.TimelineFilters) ([]model.TimelineItem, error) {
	timeline, err := s.repo.GetWithFilters(filters)
	// Todo: use requestAuthor in anonymize functions
	return s.commonGetterHandler(timeline, err, true, true)
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
		timelineItem.Participant.Challenge = nil
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
		isAnonym = s.AnonymizeItem(&(*timeline)[i]) || isAnonym
	}
	return isAnonym
}

func (s *TimelineService) AnonymizeItem(timelineItem *model.TimelineItem) bool {
	s.srvChallenges.Anonymize(&timelineItem.Challenge)
	if timelineItem.Participant != nil {
		s.srvParticipants.Anonymize(timelineItem.Participant, "")
	}

	// TODO: Anonymize for suggestion
	if timelineItem.Suggestion != nil {
		// timelineItem.Suggestion.Anonymous
	}

	isAnonym := (timelineItem.Type == "CREATE_CHALLENGE" && !timelineItem.Challenge.ShowAuthor) ||
		((timelineItem.Type == "ACCEPT_CHALLENGE" || timelineItem.Type == "ADD_NOTE") && timelineItem.Participant.Anonymous) ||
		(timelineItem.Type == "ADD_SUGGESTION" && !timelineItem.Suggestion.Anonymous)

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
