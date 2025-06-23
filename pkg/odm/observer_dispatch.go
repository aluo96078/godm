package odm

import "sort"

// observer_dispatch.go - 執行 Observer 通知流程，依照類型、事件與優先順序觸發
// Executes observer notification flows, invoking by type, event, and priority.

func (o *GODM) notifyCreating() error {
	// 處理創建階段的通知
	// Handles notifications for the creating stage.
	observers := append(globalObservers, o.Observers...)
	sort.SliceStable(observers, func(i, j int) bool {
		pi, pj := getObserverPriority(observers[i]), getObserverPriority(observers[j])
		return pi > pj
	})

	for _, observer := range observers {
		if t, ok := observer.(TypedObserver); ok && !t.Accepts(o.Model) {
			continue
		}
		if f, ok := observer.(EventFilter); ok && !f.InterestedIn("creating") {
			continue
		}
		if err := observer.Creating(o.Model); err != nil {
			if observerErrorHandler != nil {
				observerErrorHandler(err, "creating", o.Model)
			}
		}
	}
	return nil
}

func (o *GODM) notifyCreated() error {
	// 處理創建後階段的通知
	// Handles notifications for the created stage.
	observers := append(globalObservers, o.Observers...)
	sort.SliceStable(observers, func(i, j int) bool {
		pi, pj := getObserverPriority(observers[i]), getObserverPriority(observers[j])
		return pi > pj
	})

	for _, observer := range observers {
		if t, ok := observer.(TypedObserver); ok && !t.Accepts(o.Model) {
			continue
		}
		if f, ok := observer.(EventFilter); ok && !f.InterestedIn("created") {
			continue
		}
		if err := observer.Created(o.Model); err != nil {
			if observerErrorHandler != nil {
				observerErrorHandler(err, "created", o.Model)
			}
		}
	}
	return nil
}

func (o *GODM) notifyUpdating() error {
	// 處理更新階段的通知
	// Handles notifications for the updating stage.
	observers := append(globalObservers, o.Observers...)
	sort.SliceStable(observers, func(i, j int) bool {
		pi, pj := getObserverPriority(observers[i]), getObserverPriority(observers[j])
		return pi > pj
	})

	for _, observer := range observers {
		if t, ok := observer.(TypedObserver); ok && !t.Accepts(o.Model) {
			continue
		}
		if f, ok := observer.(EventFilter); ok && !f.InterestedIn("updating") {
			continue
		}
		if err := observer.Updating(o.Model); err != nil {
			if observerErrorHandler != nil {
				observerErrorHandler(err, "updating", o.Model)
			}
		}
	}
	return nil
}

func (o *GODM) notifyUpdated() error {
	// 處理更新後階段的通知
	// Handles notifications for the updated stage.
	observers := append(globalObservers, o.Observers...)
	sort.SliceStable(observers, func(i, j int) bool {
		pi, pj := getObserverPriority(observers[i]), getObserverPriority(observers[j])
		return pi > pj
	})

	for _, observer := range observers {
		if t, ok := observer.(TypedObserver); ok && !t.Accepts(o.Model) {
			continue
		}
		if f, ok := observer.(EventFilter); ok && !f.InterestedIn("updated") {
			continue
		}
		if err := observer.Updated(o.Model); err != nil {
			if observerErrorHandler != nil {
				observerErrorHandler(err, "updated", o.Model)
			}
		}
	}
	return nil
}

func (o *GODM) notifyDeleting() error {
	// 處理刪除階段的通知
	// Handles notifications for the deleting stage.
	observers := append(globalObservers, o.Observers...)
	sort.SliceStable(observers, func(i, j int) bool {
		pi, pj := getObserverPriority(observers[i]), getObserverPriority(observers[j])
		return pi > pj
	})

	for _, observer := range observers {
		if t, ok := observer.(TypedObserver); ok && !t.Accepts(o.Model) {
			continue
		}
		if f, ok := observer.(EventFilter); ok && !f.InterestedIn("deleting") {
			continue
		}
		if err := observer.Deleting(o.Model); err != nil {
			if observerErrorHandler != nil {
				observerErrorHandler(err, "deleting", o.Model)
			}
		}
	}
	return nil
}

func (o *GODM) notifyDeleted() error {
	// 處理刪除後階段的通知
	// Handles notifications for the deleted stage.
	observers := append(globalObservers, o.Observers...)
	sort.SliceStable(observers, func(i, j int) bool {
		pi, pj := getObserverPriority(observers[i]), getObserverPriority(observers[j])
		return pi > pj
	})

	for _, observer := range observers {
		if t, ok := observer.(TypedObserver); ok && !t.Accepts(o.Model) {
			continue
		}
		if f, ok := observer.(EventFilter); ok && !f.InterestedIn("deleted") {
			continue
		}
		if err := observer.Deleted(o.Model); err != nil {
			if observerErrorHandler != nil {
				observerErrorHandler(err, "deleted", o.Model)
			}
		}
	}
	return nil
}
