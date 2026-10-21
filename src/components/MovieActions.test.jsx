import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { AuthProvider } from '../auth/AuthContext';
import { putVote, addToWishlist } from '../api/backend';
import MovieActions from './MovieActions';

jest.mock('../api/backend');
beforeEach(() => localStorage.clear());

function renderActions(loggedIn = false) {
  if (loggedIn) {
    localStorage.setItem('moviehub.token', 'member-token');
    localStorage.setItem('moviehub.member', JSON.stringify({ id: 1, displayName: 'Demo' }));
  }
  return render(<AuthProvider><MemoryRouter><MovieActions movieId={1423191} /></MemoryRouter></AuthProvider>);
}

test('ยังไม่ login ต้องเห็นลิงก์เข้าสู่ระบบและไม่มีปุ่มคะแนน', () => {
  renderActions();
  expect(screen.getByRole('link', { name: 'เข้าสู่ระบบ' })).toHaveAttribute('href', '/login');
  expect(screen.queryByRole('button', { name: '7' })).not.toBeInTheDocument();
});

test('login แล้วให้คะแนน 7 ต้องส่ง movieId คะแนน และ token ก่อนแสดงผลสำเร็จ', async () => {
  const user = userEvent.setup();
  putVote.mockResolvedValue({ movieId: 1423191, score: 7 });
  renderActions(true);
  await user.click(screen.getByRole('button', { name: '7' }));
  expect(putVote).toHaveBeenCalledWith(1423191, 7, 'member-token');
  expect(await screen.findByText('บันทึกคะแนน 7 แล้ว')).toBeInTheDocument();
});

test('เพิ่ม wishlist ล้มเหลวต้องแสดง error และหัวใจยังไม่เปลี่ยนเป็นเพิ่มแล้ว', async () => {
  const user = userEvent.setup();
  addToWishlist.mockRejectedValue(new Error('เพิ่มรายการไม่สำเร็จ'));
  renderActions(true);
  await user.click(screen.getByRole('button', { name: /🤍 เพิ่มเข้ารายการที่อยากดู/ }));
  expect(await screen.findByText('เพิ่มรายการไม่สำเร็จ')).toBeInTheDocument();
  expect(screen.getByRole('button', { name: /🤍 เพิ่มเข้ารายการที่อยากดู/ })).toBeInTheDocument();
  expect(screen.queryByRole('button', { name: /❤️ อยู่ในรายการ/ })).not.toBeInTheDocument();
});
