import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import { AuthProvider } from '../auth/AuthContext';
import { register, login } from '../api/backend';
import Register from './Register';

jest.mock('../api/backend');
beforeEach(() => localStorage.clear());

function renderRegister() {
  return render(
    <AuthProvider><MemoryRouter initialEntries={['/register']}><Routes>
      <Route path="/register" element={<Register />} />
      <Route path="/" element={<h1>หน้าหลัก</h1>} />
    </Routes></MemoryRouter></AuthProvider>
  );
}

async function fillForm(user, confirmation) {
  await user.type(screen.getByPlaceholderText('ชื่อที่แสดง'), 'Demo');
  await user.type(screen.getByPlaceholderText('อีเมล'), 'demo@example.test');
  await user.type(screen.getByPlaceholderText('รหัสผ่าน (4 ตัวขึ้นไป)'), 'popcorn2569');
  await user.type(screen.getByPlaceholderText('ยืนยันรหัสผ่าน'), confirmation);
}

test('รหัสผ่านยืนยันไม่ตรงกันต้องไม่ส่งคำขอสมัครสมาชิก', async () => {
  const user = userEvent.setup();
  renderRegister();
  await fillForm(user, 'another-password');
  await user.click(screen.getByRole('button', { name: 'สมัครสมาชิก' }));
  expect(screen.getByText('รหัสผ่านทั้งสองช่องไม่ตรงกัน')).toBeInTheDocument();
  expect(register).not.toHaveBeenCalled();
});

test('รหัสผ่านยืนยันตรงกันจึงสมัคร login และแสดงหน้าหลัก', async () => {
  const user = userEvent.setup();
  register.mockResolvedValue({ id: 1 });
  login.mockResolvedValue({ token: 'signed-token', member: { id: 1, displayName: 'Demo' } });
  renderRegister();
  await fillForm(user, 'popcorn2569');
  await user.click(screen.getByRole('button', { name: 'สมัครสมาชิก' }));
  expect(await screen.findByRole('heading', { name: 'หน้าหลัก' })).toBeInTheDocument();
  expect(register).toHaveBeenCalledWith('demo@example.test', 'popcorn2569', 'Demo');
  expect(login).toHaveBeenCalledWith('demo@example.test', 'popcorn2569');
});
