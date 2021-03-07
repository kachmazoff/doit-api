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
	if err != nil {
		return []model.TimelineItem{}, err
	}

	for i := 0; i < len(timeline); i++ {
		err := s.EnrichItem(&timeline[i])
		if err != nil {
			return []model.TimelineItem{}, err
		}
	}

	for i := 0; i < len(timeline); i++ {
		s.AnonymizeItem(&timeline[i])
	}

	return timeline, nil
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
			return nil
		}

		timelineItem.Note = &note
	}

	if timelineItem.SuggestionId != nil {
		suggestion, err := s.srvSuggestions.GetById(*timelineItem.SuggestionId)

		if err != nil {
			return nil
		}

		timelineItem.Suggestion = &suggestion
	}

	return nil
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
