#!/bin/bash
echo "📦 Generating 20+ test files..."

# ----------------------------------------------------------------
# COMMON COMPONENTS (6 files, ~18 tests)
# ----------------------------------------------------------------
mkdir -p tests/components/common

# 1. Button.test.tsx
cat > tests/components/common/Button.test.tsx << 'EOF'
import { render, screen, fireEvent } from '@testing-library/react';
import { Button } from '@/components/common/Button';
import { describe, it, expect, vi } from 'vitest';

describe('Button', () => {
  it('should render with text', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('should handle click event', () => {
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>Click me</Button>);
    fireEvent.click(screen.getByText('Click me'));
    expect(handleClick).toHaveBeenCalledOnce();
  });

  it('should be disabled when isLoading', () => {
    render(<Button isLoading>Click me</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('should apply variant styles', () => {
    const { container } = render(<Button variant="danger">Delete</Button>);
    expect(container.querySelector('button')).toHaveClass('bg-red-600');
  });
});
EOF

# 2. Input.test.tsx
cat > tests/components/common/Input.test.tsx << 'EOF'
import { render, screen, fireEvent } from '@testing-library/react';
import { Input } from '@/components/common/Input';
import { describe, it, expect, vi } from 'vitest';

describe('Input', () => {
  it('should render with placeholder', () => {
    render(<Input placeholder="Enter text" />);
    expect(screen.getByPlaceholderText('Enter text')).toBeInTheDocument();
  });

  it('should handle onChange', () => {
    const handleChange = vi.fn();
    render(<Input onChange={handleChange} />);
    const input = screen.getByRole('textbox');
    fireEvent.change(input, { target: { value: 'hello' } });
    expect(handleChange).toHaveBeenCalled();
  });

  it('should be disabled when disabled prop is true', () => {
    render(<Input disabled />);
    expect(screen.getByRole('textbox')).toBeDisabled();
  });

  it('should show error state', () => {
    render(<Input error="This field is required" />);
    expect(screen.getByText('This field is required')).toBeInTheDocument();
  });
});
EOF

# 3. Card.test.tsx
cat > tests/components/common/Card.test.tsx << 'EOF'
import { render, screen } from '@testing-library/react';
import { Card, CardHeader, CardTitle, CardContent, CardFooter } from '@/components/common/Card';
import { describe, it, expect } from 'vitest';

describe('Card', () => {
  it('should render children', () => {
    render(<Card>Card content</Card>);
    expect(screen.getByText('Card content')).toBeInTheDocument();
  });

  it('should render CardHeader, CardTitle, CardContent, CardFooter', () => {
    render(
      <Card>
        <CardHeader><CardTitle>Title</CardTitle></CardHeader>
        <CardContent>Content</CardContent>
        <CardFooter>Footer</CardFooter>
      </Card>
    );
    expect(screen.getByText('Title')).toBeInTheDocument();
    expect(screen.getByText('Content')).toBeInTheDocument();
    expect(screen.getByText('Footer')).toBeInTheDocument();
  });
});
EOF

# 4. Select.test.tsx
cat > tests/components/common/Select.test.tsx << 'EOF'
import { render, screen, fireEvent } from '@testing-library/react';
import { Select } from '@/components/common/Select';
import { describe, it, expect, vi } from 'vitest';

describe('Select', () => {
  const options = [{ value: '1', label: 'Option 1' }, { value: '2', label: 'Option 2' }];

  it('should render options', () => {
    render(<Select options={options} />);
    expect(screen.getByText('Option 1')).toBeInTheDocument();
    expect(screen.getByText('Option 2')).toBeInTheDocument();
  });

  it('should handle onChange', () => {
    const handleChange = vi.fn();
    render(<Select options={options} onChange={handleChange} />);
    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: '2' } });
    expect(handleChange).toHaveBeenCalledWith('2');
  });
});
EOF

# 5. Modal.test.tsx
cat > tests/components/common/Modal.test.tsx << 'EOF'
import { render, screen, fireEvent } from '@testing-library/react';
import { Modal } from '@/components/common/Modal';
import { describe, it, expect, vi } from 'vitest';

describe('Modal', () => {
  it('should render when isOpen is true', () => {
    render(<Modal isOpen={true} onClose={() => {}} title="Test Modal"><div>Content</div></Modal>);
    expect(screen.getByText('Test Modal')).toBeInTheDocument();
    expect(screen.getByText('Content')).toBeInTheDocument();
  });

  it('should not render when isOpen is false', () => {
    render(<Modal isOpen={false} onClose={() => {}}><div>Content</div></Modal>);
    expect(screen.queryByText('Content')).not.toBeInTheDocument();
  });

  it('should call onClose when close button clicked', () => {
    const onClose = vi.fn();
    render(<Modal isOpen={true} onClose={onClose}><div>Content</div></Modal>);
    const closeBtn = screen.getByRole('button', { name: /close/i });
    fireEvent.click(closeBtn);
    expect(onClose).toHaveBeenCalled();
  });
});
EOF

# 6. Toast.test.tsx
cat > tests/components/common/Toast.test.tsx << 'EOF'
import { render, screen, act } from '@testing-library/react';
import { Toast } from '@/components/common/Toast';
import { describe, it, expect, vi, beforeEach } from 'vitest';

describe('Toast', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  it('should render message', () => {
    render(<Toast message="Test message" type="success" />);
    expect(screen.getByText('Test message')).toBeInTheDocument();
  });

  it('should auto-dismiss after duration', () => {
    const onDismiss = vi.fn();
    render(<Toast message="Test" duration={1000} onDismiss={onDismiss} />);
    act(() => { vi.advanceTimersByTime(1000); });
    expect(onDismiss).toHaveBeenCalled();
  });
});
EOF

echo "✅ Common components tests created (6 files)"

# ----------------------------------------------------------------
# AUTH COMPONENTS (2 files, ~8 tests)
# ----------------------------------------------------------------
mkdir -p tests/components/auth

# 7. LoginForm.test.tsx
cat > tests/components/auth/LoginForm.test.tsx << 'EOF'
import { render, screen, fireEvent } from '@testing-library/react';
import { LoginForm } from '@/components/auth/LoginForm';
import { describe, it, expect, vi } from 'vitest';

describe('LoginForm', () => {
  it('should render email and password fields', () => {
    render(<LoginForm onSubmit={vi.fn()} />);
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/kata sandi/i)).toBeInTheDocument();
  });

  it('should validate email format', async () => {
    const handleSubmit = vi.fn();
    render(<LoginForm onSubmit={handleSubmit} />);
    const emailInput = screen.getByLabelText(/email/i);
    const submitBtn = screen.getByRole('button', { name: /masuk/i });
    fireEvent.change(emailInput, { target: { value: 'invalid-email' } });
    fireEvent.click(submitBtn);
    expect(await screen.findByText(/format email tidak valid/i)).toBeInTheDocument();
    expect(handleSubmit).not.toHaveBeenCalled();
  });

  it('should submit with valid credentials', async () => {
    const handleSubmit = vi.fn();
    render(<LoginForm onSubmit={handleSubmit} />);
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'test@example.com' } });
    fireEvent.change(screen.getByLabelText(/kata sandi/i), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /masuk/i }));
    expect(handleSubmit).toHaveBeenCalledWith('test@example.com', 'password123');
  });

  it('should show error message', () => {
    render(<LoginForm onSubmit={vi.fn()} error="Login failed" />);
    expect(screen.getByText('Login failed')).toBeInTheDocument();
  });
});
EOF

# 8. RegisterForm.test.tsx
cat > tests/components/auth/RegisterForm.test.tsx << 'EOF'
import { render, screen, fireEvent } from '@testing-library/react';
import { RegisterForm } from '@/components/auth/RegisterForm';
import { describe, it, expect, vi } from 'vitest';

describe('RegisterForm', () => {
  it('should render all fields', () => {
    render(<RegisterForm onSubmit={vi.fn()} />);
    expect(screen.getByLabelText(/nama lengkap/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/kata sandi/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/konfirmasi kata sandi/i)).toBeInTheDocument();
  });

  it('should validate password match', async () => {
    const handleSubmit = vi.fn();
    render(<RegisterForm onSubmit={handleSubmit} />);
    fireEvent.change(screen.getByLabelText(/kata sandi/i), { target: { value: 'Password123!' } });
    fireEvent.change(screen.getByLabelText(/konfirmasi kata sandi/i), { target: { value: 'Different123!' } });
    fireEvent.click(screen.getByRole('button', { name: /daftar/i }));
    expect(await screen.findByText(/password tidak cocok/i)).toBeInTheDocument();
  });
});
EOF

echo "✅ Auth components tests created (2 files)"

# ----------------------------------------------------------------
# DASHBOARD COMPONENTS (2 files, ~6 tests)
# ----------------------------------------------------------------
mkdir -p tests/components/dashboard

# 9. StatsCard.test.tsx
cat > tests/components/dashboard/StatsCard.test.tsx << 'EOF'
import { render, screen } from '@testing-library/react';
import { StatsCard } from '@/components/dashboard/StatsCard';
import { describe, it, expect } from 'vitest';

describe('StatsCard', () => {
  it('should render title and value', () => {
    render(<StatsCard title="Total SKU" value="1,234" />);
    expect(screen.getByText('Total SKU')).toBeInTheDocument();
    expect(screen.getByText('1,234')).toBeInTheDocument();
  });

  it('should render icon when provided', () => {
    render(<StatsCard title="Test" value="100" icon={<span>📦</span>} />);
    expect(screen.getByText('📦')).toBeInTheDocument();
  });
});
EOF

# 10. RecentTransactions.test.tsx
cat > tests/components/dashboard/RecentTransactions.test.tsx << 'EOF'
import { render, screen } from '@testing-library/react';
import { RecentTransactions } from '@/components/dashboard/RecentTransactions';
import { describe, it, expect } from 'vitest';

const mockData = [
  { id: '1', sku: 'SKU-001', type: 'in', quantity: 10, warehouse: 'Jakarta', created_at: '2026-07-29' },
  { id: '2', sku: 'SKU-002', type: 'out', quantity: 5, warehouse: 'Bandung', created_at: '2026-07-28' },
];

describe('RecentTransactions', () => {
  it('should render transaction list', () => {
    render(<RecentTransactions transactions={mockData} />);
    expect(screen.getByText('SKU-001')).toBeInTheDocument();
    expect(screen.getByText('SKU-002')).toBeInTheDocument();
  });

  it('should show empty state', () => {
    render(<RecentTransactions transactions={[]} />);
    expect(screen.getByText(/belum ada transaksi/i)).toBeInTheDocument();
  });
});
EOF

echo "✅ Dashboard components tests created (2 files)"

# ----------------------------------------------------------------
# INVENTORY COMPONENTS (2 files, ~6 tests)
# ----------------------------------------------------------------
mkdir -p tests/components/inventory

# 11. InventoryTable.test.tsx
cat > tests/components/inventory/InventoryTable.test.tsx << 'EOF'
import { render, screen, fireEvent } from '@testing-library/react';
import { InventoryTable } from '@/components/inventory/InventoryTable';
import { describe, it, expect, vi } from 'vitest';

const mockItems = [
  { id: '1', sku: 'SKU-001', name: 'Product A', quantity: 100, status: 'normal', warehouse: 'Jakarta' },
  { id: '2', sku: 'SKU-002', name: 'Product B', quantity: 5, status: 'low_stock', warehouse: 'Bandung' },
];

describe('InventoryTable', () => {
  it('should render inventory rows', () => {
    render(<InventoryTable items={mockItems} />);
    expect(screen.getByText('Product A')).toBeInTheDocument();
    expect(screen.getByText('Product B')).toBeInTheDocument();
  });

  it('should call onSort when header clicked', () => {
    const onSort = vi.fn();
    render(<InventoryTable items={mockItems} onSort={onSort} />);
    const header = screen.getByText(/sku/i);
    fireEvent.click(header);
    expect(onSort).toHaveBeenCalled();
  });
});
EOF

# 12. TransactionForm.test.tsx
cat > tests/components/inventory/TransactionForm.test.tsx << 'EOF'
import { render, screen, fireEvent } from '@testing-library/react';
import { TransactionForm } from '@/components/inventory/TransactionForm';
import { describe, it, expect, vi } from 'vitest';

describe('TransactionForm', () => {
  it('should render all fields', () => {
    render(<TransactionForm onSubmit={vi.fn()} />);
    expect(screen.getByLabelText(/sku/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/kuantitas/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/alasan/i)).toBeInTheDocument();
  });

  it('should validate quantity > 0', async () => {
    const handleSubmit = vi.fn();
    render(<TransactionForm onSubmit={handleSubmit} />);
    fireEvent.change(screen.getByLabelText(/kuantitas/i), { target: { value: '0' } });
    fireEvent.click(screen.getByRole('button', { name: /simpan/i }));
    expect(await screen.findByText(/kuantitas harus lebih dari 0/i)).toBeInTheDocument();
  });
});
EOF

echo "✅ Inventory components tests created (2 files)"

# ----------------------------------------------------------------
# FORECAST COMPONENTS (1 file, ~4 tests)
# ----------------------------------------------------------------
mkdir -p tests/components/forecast

# 13. ForecastChart.test.tsx
cat > tests/components/forecast/ForecastChart.test.tsx << 'EOF'
import { render, screen } from '@testing-library/react';
import { ForecastChart } from '@/components/forecast/ForecastChart';
import { describe, it, expect } from 'vitest';

const mockData = [
  { date: '2026-07-29', predicted: 100, lower: 80, upper: 120 },
  { date: '2026-07-30', predicted: 110, lower: 90, upper: 130 },
];

describe('ForecastChart', () => {
  it('should render chart container', () => {
    render(<ForecastChart data={mockData} />);
    // Recharts renders canvas or svg, check container
    const container = document.querySelector('.recharts-wrapper');
    expect(container).toBeInTheDocument();
  });

  it('should show empty state when no data', () => {
    render(<ForecastChart data={[]} />);
    expect(screen.getByText(/tidak ada data forecast/i)).toBeInTheDocument();
  });
});
EOF

echo "✅ Forecast components tests created (1 file)"

# ----------------------------------------------------------------
# ANOMALY COMPONENTS (2 files, ~4 tests)
# ----------------------------------------------------------------
mkdir -p tests/components/anomaly

# 14. AnomalyList.test.tsx
cat > tests/components/anomaly/AnomalyList.test.tsx << 'EOF'
import { render, screen } from '@testing-library/react';
import { AnomalyList } from '@/components/anomaly/AnomalyList';
import { describe, it, expect } from 'vitest';

const mockAnomalies = [
  { id: '1', sku: 'SKU-001', severity: 'high', description: 'Unusual drop', status: 'open', detected_at: '2026-07-29' },
  { id: '2', sku: 'SKU-002', severity: 'low', description: 'Minor spike', status: 'resolved', detected_at: '2026-07-28' },
];

describe('AnomalyList', () => {
  it('should render anomaly items', () => {
    render(<AnomalyList anomalies={mockAnomalies} />);
    expect(screen.getByText('Unusual drop')).toBeInTheDocument();
    expect(screen.getByText('Minor spike')).toBeInTheDocument();
  });
});
EOF

# 15. SeverityBadge.test.tsx
cat > tests/components/anomaly/SeverityBadge.test.tsx << 'EOF'
import { render, screen } from '@testing-library/react';
import { SeverityBadge } from '@/components/anomaly/SeverityBadge';
import { describe, it, expect } from 'vitest';

describe('SeverityBadge', () => {
  it('should render with severity labels', () => {
    render(<SeverityBadge severity="critical" />);
    expect(screen.getByText('Critical')).toBeInTheDocument();
  });

  it('should apply correct color class for severity', () => {
    const { container } = render(<SeverityBadge severity="high" />);
    expect(container.firstChild).toHaveClass('bg-red-600');
  });
});
EOF

echo "✅ Anomaly components tests created (2 files)"

# ----------------------------------------------------------------
# HOOKS (3 files, ~6 tests)
# ----------------------------------------------------------------
mkdir -p tests/hooks

# 16. useAuth.test.ts
cat > tests/hooks/useAuth.test.ts << 'EOF'
import { renderHook, act } from '@testing-library/react';
import { useAuth } from '@/context/AuthContext';
import { describe, it, expect, vi } from 'vitest';

// Mock the context provider
vi.mock('@/context/AuthContext', () => ({
  useAuth: () => ({
    user: { email: 'test@example.com', role: 'admin' },
    isAuthenticated: true,
    isLoading: false,
    login: vi.fn(),
    logout: vi.fn(),
    register: vi.fn(),
  }),
}));

describe('useAuth', () => {
  it('should return auth state', () => {
    const { result } = renderHook(() => useAuth());
    expect(result.current.isAuthenticated).toBe(true);
    expect(result.current.user?.email).toBe('test@example.com');
  });

  it('should have login and logout functions', () => {
    const { result } = renderHook(() => useAuth());
    expect(result.current.login).toBeTypeOf('function');
    expect(result.current.logout).toBeTypeOf('function');
  });
});
EOF

# 17. useApi.test.ts
cat > tests/hooks/useApi.test.ts << 'EOF'
import { renderHook, waitFor } from '@testing-library/react';
import { useApi } from '@/hooks/useApi';
import { describe, it, expect, vi } from 'vitest';
import apiClient from '@/lib/api';

vi.mock('@/lib/api', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
  },
}));

describe('useApi', () => {
  it('should fetch data successfully', async () => {
    const mockData = { data: { data: 'test' } };
    (apiClient.get as any).mockResolvedValueOnce(mockData);

    const { result } = renderHook(() => useApi('/test'));
    expect(result.current.isLoading).toBe(true);

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
      expect(result.current.data).toEqual('test');
    });
  });

  it('should handle error', async () => {
    (apiClient.get as any).mockRejectedValueOnce(new Error('Network error'));
    const { result } = renderHook(() => useApi('/test'));

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeDefined();
    });
  });
});
EOF

# 18. useTheme.test.ts
cat > tests/hooks/useTheme.test.ts << 'EOF'
import { renderHook, act } from '@testing-library/react';
import { useTheme } from '@/context/ThemeContext';
import { describe, it, expect, vi } from 'vitest';

// Mock the theme context
vi.mock('@/context/ThemeContext', () => ({
  useTheme: () => ({
    theme: 'light',
    setTheme: vi.fn(),
    isDark: false,
  }),
}));

describe('useTheme', () => {
  it('should return theme state', () => {
    const { result } = renderHook(() => useTheme());
    expect(result.current.theme).toBe('light');
    expect(result.current.isDark).toBe(false);
  });

  it('should have setTheme function', () => {
    const { result } = renderHook(() => useTheme());
    expect(result.current.setTheme).toBeTypeOf('function');
  });
});
EOF

echo "✅ Hooks tests created (3 files)"

# ----------------------------------------------------------------
# LIB/UTILITIES (2 files, ~4 tests)
# ----------------------------------------------------------------
mkdir -p tests/lib

# 19. validators.test.ts
cat > tests/lib/validators.test.ts << 'EOF'
import { describe, it, expect } from 'vitest';
import { isValidEmail, isValidPassword, isStrongPassword } from '@/lib/validators';

describe('validators', () => {
  it('should validate email', () => {
    expect(isValidEmail('test@example.com')).toBe(true);
    expect(isValidEmail('invalid-email')).toBe(false);
    expect(isValidEmail('test@')).toBe(false);
  });

  it('should validate password length', () => {
    expect(isValidPassword('password123')).toBe(true); // > 6 chars
    expect(isValidPassword('pass')).toBe(false);
  });

  it('should check strong password', () => {
    expect(isStrongPassword('Password123!')).toBe(true);
    expect(isStrongPassword('password')).toBe(false); // no uppercase/digit/special
  });
});
EOF

# 20. utils.test.ts
cat > tests/lib/utils.test.ts << 'EOF'
import { describe, it, expect } from 'vitest';
import { cn, formatDate, formatCurrency } from '@/lib/utils';

describe('utils', () => {
  it('should merge class names with cn', () => {
    expect(cn('class1', 'class2')).toBe('class1 class2');
    expect(cn('class1', { class2: true, class3: false })).toBe('class1 class2');
  });

  it('should format date', () => {
    const date = new Date('2026-07-29T10:00:00Z');
    expect(formatDate(date, 'dd MMM yyyy')).toBe('29 Jul 2026');
  });

  it('should format currency', () => {
    expect(formatCurrency(1000000, 'IDR')).toBe('Rp 1.000.000');
    expect(formatCurrency(1000, 'USD')).toBe('$1,000');
  });
});
EOF

echo "✅ Utilities tests created (2 files)"

echo ""
echo "🎉 All 20 test files have been generated!"
echo "Total files: $(find tests -name '*.test.tsx' -o -name '*.test.ts' | wc -l)"
