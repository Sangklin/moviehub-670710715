import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from '../auth/AuthContext';
import { login } from '../api/backend';
import Login from './Login';

jest.mock('../api/backend');
beforeEach(() => localStorage.clear());

function renderLogin() {
  return render(
    <AuthProvider>
      <MemoryRouter initialEntries={[{ pathname: '/login', state: { from: '/me/wishlist' } }]}>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/me/wishlist" element={<h1>รายการส่วนตัว</h1>} />
        </Routes>
      </MemoryRouter>
    </AuthProvider>
  );
}

test('รหัสผ่านผิดต้องแสดงข้อความจาก server และให้ส่งฟอร์มใหม่ได้', async () => {
  const user = userEvent.setup();
  login.mockRejectedValue(new Error('อีเมลหรือรหัสผ่านไม่ถูกต้อง'));
  renderLogin();
  await user.click(screen.getByRole('button', { name: 'เข้าสู่ระบบ' }));
  expect(await screen.findByText('อีเมลหรือรหัสผ่านไม่ถูกต้อง')).toBeInTheDocument();
  expect(screen.getByRole('button', { name: 'เข้าสู่ระบบ' })).toBeEnabled();
  expect(localStorage.getItem('moviehub.token')).toBeNull();
});

test('login สำเร็จต้องจำ token และกลับไปหน้าที่ต้องการก่อนเข้าสู่ระบบ', async () => {
  const user = userEvent.setup();
  login.mockResolvedValue({ token: 'signed-token', member: { id: 1, displayName: 'Demo' } });
  renderLogin();
  await user.click(screen.getByRole('button', { name: 'เข้าสู่ระบบ' }));
  expect(await screen.findByRole('heading', { name: 'รายการส่วนตัว' })).toBeInTheDocument();
  expect(localStorage.getItem('moviehub.token')).toBe('signed-token');
});
