package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cloudfoundry-incubator/bbs/db"
	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/pivotal-golang/lager"
)

type TaskHandler struct {
	db     db.TaskDB
	logger lager.Logger
}

func NewTaskHandler(logger lager.Logger, db db.TaskDB) *TaskHandler {
	return &TaskHandler{
		db:     db,
		logger: logger.Session("task-handler"),
	}
}

func (h *TaskHandler) Tasks(w http.ResponseWriter, req *http.Request) {
	domain := req.FormValue("domain")
	cellID := req.FormValue("cell_id")
	logger := h.logger.Session("tasks", lager.Data{
		"domain":  domain,
		"cell_id": cellID,
	})

	if domain != "" && cellID != "" {
		writeBadRequestResponse(w, models.InvalidRequest, errors.New("too many filters"))
		return
	}

	tasks, err := h.db.Tasks(h.logger, taskFilter(domain, cellID))
	if err != nil {
		logger.Error("failed-to-fetch-tasks", err)
		writeInternalServerErrorResponse(w, err)
		return
	}

	writeProtoResponse(w, http.StatusOK, tasks)
}

func taskFilter(domain string, cellID string) db.TaskFilter {
	if domain != "" {
		return func(t *models.Task) bool {
			return domain == t.Domain
		}
	}
	if cellID != "" {
		return func(t *models.Task) bool {
			return cellID == t.CellId
		}
	}

	return nil
}

func (h *TaskHandler) TaskByGuid(w http.ResponseWriter, req *http.Request) {
	taskGuid := req.FormValue(":task_guid")
	logger := h.logger.Session("task-by-guid", lager.Data{
		"task_guid": taskGuid,
	})

	task, err := h.db.TaskByGuid(h.logger, taskGuid)
	if err == models.ErrResourceNotFound {
		writeNotFoundResponse(w, err)
		return
	}
	if err != nil {
		logger.Error("failed-to-fetch-task", err)
		writeInternalServerErrorResponse(w, err)
		return
	}

	writeProtoResponse(w, http.StatusOK, task)
}

func (h *TaskHandler) DesireTask(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.Session("desire-task")

	request := &models.DesireTaskRequest{}
	response := &models.TaskLifecycleResponse{}

	response.Error = parseRequest(logger, req, request)
	if response.Error == nil {
		response.Error = h.db.DesireTask(h.logger, request.TaskDefinition, request.TaskGuid, request.Domain)
	}

	writeResponse(w, response)
}

func (h *TaskHandler) StartTask(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.Session("start-task")

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("failed-to-read-body", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to read request body: %s", err))
		return
	}

	startReq := &models.StartTaskRequest{}
	err = startReq.Unmarshal(data)
	if err != nil {
		logger.Error("failed-to-unmarshal", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to unmarshal: %s", err))
		return
	}

	if err := startReq.Validate(); err != nil {
		logger.Error("invalid-request", err)
		writeBadRequestResponse(w, models.InvalidRequest, err)
		return
	}

	taskGuid := startReq.TaskGuid
	cellID := startReq.CellId
	logger = logger.WithData(lager.Data{"task-guid": taskGuid, "cell-id": cellID})

	shouldStart, startErr := h.db.StartTask(logger, startReq.TaskGuid, startReq.CellId)
	if startErr != nil {
		logger.Error("failed-to-start-task", startErr)
		writeInternalServerErrorResponse(w, startErr)
		return
	}
	logger.Info("succeeded-start-task")
	writeProtoResponse(w, http.StatusOK, &models.StartTaskResponse{ShouldStart: shouldStart})
}

func (h *TaskHandler) CancelTask(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.Session("cancel-task")

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("failed-to-read-body", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to read request body: %s", err))
		return
	}

	request := &models.TaskGuidRequest{}
	err = request.Unmarshal(data)
	if err != nil {
		logger.Error("failed-to-unmarshal-task", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to unmarshal cancel task request: %s", err))
		return
	}
	logger.Debug("parsed-request-body", lager.Data{"request": request})
	// if err := request.Validate(); err != nil {
	// 	logger.Error("invalid-request", err)
	// 	writeBadRequestResponse(w, models.InvalidRequest, err)
	// 	return
	// }

	modelErr := h.db.CancelTask(logger, request.TaskGuid)
	if modelErr != nil {
		logger.Error("failed-to-cancel-task", modelErr)
		if modelErr.Type == models.InvalidRecord {
			writeBadRequestResponse(w, models.InvalidRecord, modelErr)
		} else {
			writeInternalServerErrorResponse(w, modelErr)
		}
		return
	}

	writeEmptyResponse(w, http.StatusNoContent)
}

func (h *TaskHandler) FailTask(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.Session("fail-task")

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("failed-to-read-body", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to read request body: %s", err))
		return
	}

	dbReq := &models.FailTaskRequest{}
	err = dbReq.Unmarshal(data)
	if err != nil {
		logger.Error("failed-to-unmarshal", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to unmarshal: %s", err))
		return
	}
	// if err := dbReq.Validate(); err != nil {
	// 	logger.Error("invalid-request", err)
	// 	writeBadRequestResponse(w, models.InvalidRequest, err)
	// 	return
	// }

	logger.Debug("parsed-request-body", lager.Data{"request": dbReq})

	modelErr := h.db.FailTask(logger, dbReq.TaskGuid, dbReq.FailureReason)
	if modelErr != nil {
		logger.Error("failed-to-fail-task", modelErr)
		writeInternalServerErrorResponse(w, modelErr)
		return
	}
	logger.Info("succeeded-fail-task")
	writeEmptyResponse(w, http.StatusNoContent)
}

func (h *TaskHandler) CompleteTask(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.Session("complete-task")

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("failed-to-read-body", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to read request body: %s", err))
		return
	}

	dbReq := &models.CompleteTaskRequest{}
	err = dbReq.Unmarshal(data)
	if err != nil {
		logger.Error("failed-to-unmarshal", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to unmarshal: %s", err))
		return
	}
	// if err := dbReq.Validate(); err != nil {
	// 	logger.Error("invalid-request", err)
	// 	writeBadRequestResponse(w, models.InvalidRequest, err)
	// 	return
	// }

	logger.Debug("parsed-request-body", lager.Data{"request": dbReq})

	modelErr := h.db.CompleteTask(logger, dbReq.TaskGuid, dbReq.CellId, dbReq.Failed, dbReq.FailureReason, dbReq.Result)
	if modelErr != nil {
		logger.Error("failed-to-complete-task", modelErr)
		writeInternalServerErrorResponse(w, modelErr)
		return
	}
	logger.Info("succeeded-complete-task")
	writeEmptyResponse(w, http.StatusNoContent)
}

func (h *TaskHandler) ResolvingTask(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.Session("resolving-task")

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("failed-to-read-body", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to read request body: %s", err))
		return
	}

	request := &models.TaskGuidRequest{}
	err = request.Unmarshal(data)
	if err != nil {
		logger.Error("failed-to-unmarshal-task", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to unmarshal reesolving task request: %s", err))
		return
	}
	logger.Debug("parsed-request-body", lager.Data{"request": request})
	// if err := request.Validate(); err != nil {
	// 	logger.Error("invalid-request", err)
	// 	writeBadRequestResponse(w, models.InvalidRequest, err)
	// 	return
	// }

	modelErr := h.db.ResolvingTask(logger, request.TaskGuid)
	if modelErr != nil {
		logger.Error("failed-resolving-task", modelErr)
		if modelErr.Type == models.InvalidRecord {
			writeBadRequestResponse(w, models.InvalidRecord, modelErr)
		} else {
			writeInternalServerErrorResponse(w, modelErr)
		}
		return
	}

	writeEmptyResponse(w, http.StatusNoContent)
}

func (h *TaskHandler) ResolveTask(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.Session("resolve-task")

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("failed-to-read-body", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to read request body: %s", err))
		return
	}

	request := &models.TaskGuidRequest{}
	err = request.Unmarshal(data)
	if err != nil {
		logger.Error("failed-to-unmarshal-task", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to unmarshal resolve task request: %s", err))
		return
	}
	logger.Debug("parsed-request-body", lager.Data{"request": request})
	// if err := request.Validate(); err != nil {
	// 	logger.Error("invalid-request", err)
	// 	writeBadRequestResponse(w, models.InvalidRequest, err)
	// 	return
	// }

	modelErr := h.db.ResolveTask(logger, request.TaskGuid)
	if modelErr != nil {
		logger.Error("failed-to-resolve-task", modelErr)
		if modelErr.Type == models.InvalidRecord {
			writeBadRequestResponse(w, models.InvalidRecord, modelErr)
		} else {
			writeInternalServerErrorResponse(w, modelErr)
		}
		return
	}

	writeEmptyResponse(w, http.StatusNoContent)
}

func (h *TaskHandler) ConvergeTasks(w http.ResponseWriter, req *http.Request) {
	logger := h.logger.Session("converge-tasks")

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("failed-to-read-body", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to read request body: %s", err))
		return
	}

	request := &models.ConvergeTasksRequest{}
	err = request.Unmarshal(data)
	if err != nil {
		logger.Error("failed-to-unmarshal-task", err)
		writeBadRequestResponse(w, models.InvalidRequest, fmt.Errorf("failed to unmarshal converge tasks request: %s", err))
		return
	}
	logger.Debug("parsed-request-body", lager.Data{"request": request})
	// if err := request.Validate(); err != nil {
	// 	logger.Error("invalid-request", err)
	// 	writeBadRequestResponse(w, models.InvalidRequest, err)
	// 	return
	// }

	h.db.ConvergeTasks(
		logger,
		time.Duration(request.KickTaskDuration),
		time.Duration(request.ExpirePendingTaskDuration),
		time.Duration(request.ExpireCompletedTaskDuration),
	)
	writeEmptyResponse(w, http.StatusNoContent)
}
