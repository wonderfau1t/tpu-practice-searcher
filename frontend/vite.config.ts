import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';

export default defineConfig({
  plugins: [react()],
  server: {
    host: true, // Разрешаем доступ извне
    port: 5173, // Убедись, что совпадает с ngrok
    strictPort: true, // Использовать только этот порт
    allowedHosts: ['.ngrok-free.app'], // Разрешаем ngrok-домены
  },
})
