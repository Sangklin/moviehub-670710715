import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { AuthProvider } from '../auth/AuthContext';
import { getWishlist } from '../api/backend';
import Wishlist from './Wishlist';

jest.mock('../api/backend');

beforeEach(() => {
  localStorage.clear();
  localStorage.setItem('moviehub.token', 'member-token');
  localStorage.setItem('moviehub.member', JSON.stringify({ id: 1, displayName: 'Demo' }));
});

function renderWishlist() {
  return render(<AuthProvider><MemoryRouter><Wishlist /></MemoryRouter></AuthProvider>);
}

test('โหลดรายการอยากดูด้วย token ของสมาชิกและแสดงหนังที่ server คืนมา', async () => {
  getWishlist.mockResolvedValue({ items: [{ id: 1, title: 'Parasite', poster: null }] });
  renderWishlist();
  expect(await screen.findByText('Parasite')).toBeInTheDocument();
  expect(getWishlist).toHaveBeenCalledWith('member-token');
});

test('รายการอยากดูว่างต้องแสดงข้อความไม่พบหนัง', async () => {
  getWishlist.mockResolvedValue({ items: [] });
  renderWishlist();
  expect(await screen.findByText(/ไม่พบหนัง/)).toBeInTheDocument();
});
