package http

import (
	"context"
	"errors"
	"fmt"

	"workshop/internal/domain"
	"workshop/internal/usecase"
)

func requireOrderAccess(ctx context.Context, orderID string, identity caller, assignment usecase.AssignmentUseCase, technician usecase.TechnicianUseCase) error {
	if identity.Role == domain.RoleAdministrator {
		return nil
	}
	if identity.Role != domain.RoleTechnician {
		return fmt.Errorf("%w: this service order is restricted", domain.ErrForbidden)
	}
	profile, err := technician.FindByUserID(ctx, identity.UserID)
	if err != nil {
		return fmt.Errorf("%w: the signed-in user is not a technician", domain.ErrForbidden)
	}
	active, err := assignment.FindActive(ctx, orderID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fmt.Errorf("%w: the service order has no assigned technician", domain.ErrForbidden)
		}
		return err
	}
	if active.TechnicianID != profile.ID {
		return fmt.Errorf("%w: only the assigned technician can access this service order", domain.ErrForbidden)
	}
	return nil
}
