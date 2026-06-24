import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'Sregep Dashboard',
  description: 'Finance and Pomodoro dashboard'
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
