import { Navigate } from 'react-router-dom';
import type { ReactNode } from 'react';
import { useSession } from '../shared/SessionContext';

export function ProtectedRoute({ children, administratorOnly = false }: { children: ReactNode; administratorOnly?: boolean }) {
  const { session, isAdministrator } = useSession();
  if (!session) return <Navigate to="/login" replace />;
  if (administratorOnly && !isAdministrator) return <Navigate to="/service-orders" replace />;
  return <>{children}</>;
}
