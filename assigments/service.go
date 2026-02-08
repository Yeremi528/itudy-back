package assignments

import (
	"context"
	"fmt"
	"sort"
	"time"
)

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
}

func (s *service) AssignmentsAvailables(ctx context.Context, tech, level string) ([]AssignmentTest, error) {
	// 1. Obtener Trabajadores
	workers, err := s.repository.QueryWorkersByTechAndLevel(ctx, tech, level)
	if err != nil {
		return nil, fmt.Errorf("assignments.AssignmentsAvailables: %w", err)
	}

	// MAPA CLAVE: Aquí guardaremos los horarios únicos.
	// La clave es el string de la fecha (ej: "2026-02-09T12:00:00Z")
	// Si la clave ya existe, IGNORAMOS al siguiente trabajador para ese horario.
	uniqueSlotsMap := make(map[string]AssignmentTest)

	const workStartHour = 9
	const slotDuration = 30 * time.Minute

	// Iteramos trabajadores.
	// NOTA: El orden de este array define quién tiene prioridad para tomar el horario.
	for _, worker := range workers {

		// A. Obtener agenda ocupada de ESTE trabajador
		occupiedAssignments, err := s.repository.AssignmentsByWorker(ctx, worker.ID)
		if err != nil {
			return nil, fmt.Errorf("error fetching assignments: %w", err)
		}

		// Mapa rápido para consultar si ESTE trabajador está ocupado
		workerBusyMap := make(map[string]bool)
		for _, occ := range occupiedAssignments {
			isoDate := occ.FechaAsignacion.UTC().Format(time.RFC3339)
			workerBusyMap[isoDate] = true
		}

		// B. Configurar Zona Horaria del Trabajador
		loc, err := time.LoadLocation(worker.DisponibilidadBase.ZonaHoraria) //
		if err != nil {
			loc = time.UTC
		}

		// Calcular "Hoy" a las 00:00 en su zona horaria
		nowInLoc := time.Now().In(loc)
		today := time.Date(nowInLoc.Year(), nowInLoc.Month(), nowInLoc.Day(), 0, 0, 0, 0, loc)

		// C. Loop de 7 días
		for i := 0; i < 7; i++ {
			currentDay := today.AddDate(0, 0, i)

			// Validar Día Laboral (1=Lunes ... 7=Domingo)
			weekdayGo := int(currentDay.Weekday())
			if weekdayGo == 0 {
				weekdayGo = 7
			}

			if !contains(worker.DisponibilidadBase.DiasLaborales, weekdayGo) {
				continue // Este trabajador no trabaja hoy, pasamos al siguiente día
			}

			// Calcular bloques diarios (ej: 480 min / 30 = 16 bloques)
			totalSlots := worker.DisponibilidadBase.MinutosDiarios / 30 //

			// Hora de inicio (asumimos 9 AM local, ajusta si tienes un campo 'hora_inicio')
			startTime := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(), workStartHour, 0, 0, 0, loc)

			for b := 0; b < totalSlots; b++ {
				slotTime := startTime.Add(time.Duration(b) * slotDuration)

				// CLAVE ÚNICA (UTC)
				// Usamos RFC3339 ("2026-02-08T10:00:00Z") para coincidir con el formato de Mongo
				slotISO := slotTime.UTC().Format(time.RFC3339)

				// --- LÓGICA DE DEDUPLICACIÓN ---

				// 1. ¿El hueco GLOBAL ya está cubierto por un trabajador anterior?
				if _, alreadyCovered := uniqueSlotsMap[slotISO]; alreadyCovered {
					continue // Ya tenemos a alguien para esta hora, saltar.
				}

				// 2. Si el hueco está libre, ¿ESTE trabajador está libre?
				if _, isWorkerBusy := workerBusyMap[slotISO]; !isWorkerBusy {

					// ¡Bingo! Encontramos un trabajador libre para un horario que nadie más ha tomado.
					uniqueSlotsMap[slotISO] = AssignmentTest{
						WorkerID:        worker.ID,
						TestID:          "",
						FechaAsignacion: slotTime.UTC(),
						Estado:          "disponible",
					}
				}
			}
		}
	}

	// D. Convertir Mapa a Slice ordenado (para devolver JSON limpio)
	var finalResult []AssignmentTest
	for _, slot := range uniqueSlotsMap {
		finalResult = append(finalResult, slot)
	}

	// Ordenar cronológicamente
	sort.Slice(finalResult, func(i, j int) bool {
		return finalResult[i].FechaAsignacion.Before(finalResult[j].FechaAsignacion)
	})

	return finalResult, nil
}

func (s *service) AssignmentTestByUserID(ctx context.Context, userID string) ([]AssignmentTest, error) {
	AssignmentTestByUserID, err := s.repository.AssignmentTestByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("assignments.AssignmentTestByUserID: %w", err)
	}

	return AssignmentTestByUserID, nil

}

func (s *service) CreateAssignmentsByUserID(ctx context.Context, assignment AssignmentTest) error {
	err := s.repository.CreateAssignment(ctx, assignment)
	if err != nil {
		return fmt.Errorf("assignments.CreateAssignmentsByUserID: %w", err)
	}
	return nil
}

//////////////////////////////////////////////////////////////
//////////////////////// HELPERS ////////////////////////
//////////////////////////////////////////////////////////////

func contains(slice []int, item int) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
