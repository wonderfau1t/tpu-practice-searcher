import { create } from "zustand";
import { API_URL } from "../api/instance.ts";
import { AuthResponse } from "../types/auth.ts";
import axios from "axios";

interface AuthState {
  role: 'student' | 'moderator' | 'headHR' | 'HR' | 'guest' | 'admin';
  errorMessage: string;
  isLoading: boolean;
  error: string | null;
  isAuthChecked: boolean;
  isForbidden: boolean;
  setRole: (role: 'student' | 'moderator' | 'headHR' | 'guest' | 'admin') => void;
  verifyUser: (initData: unknown) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  role: 'guest',
  errorMessage: '',
  isLoading: false,
  error: null,
  isAuthChecked: false,
  isForbidden: false,
  setRole: (role) => set({ role }),
  verifyUser: async (initData: unknown) => {
    try {
      set({ isLoading: true, isAuthChecked: false, error: null });
      const response = await axios.get<AuthResponse>(`${API_URL}/auth`, {
        headers: {
          Authorization: `tma ${initData}`,
        }
      });
      console.log('verifyUser: Response received', response.data);
      localStorage.setItem('accessToken', response.data.result.accessToken);
      set({ role: response.data.result.role });
    } catch (error: any) {
      if (error?.response?.status === 403) {
        set({
          isForbidden: true,
          errorMessage: 'Заявка на регистрацию Вашей компании находится на рассмотрении. \nО ее принятии / отклонении Вас уведомит чат-бот приложения.',
        });
        return; // не продолжаем выполнение
      }

      // Обработка других ошибок
      set({ error: error });
    } finally {
      set({ isLoading: false, isAuthChecked: true });
    }
  },
  logout: () => {
    localStorage.removeItem('accessToken');
    set({ role: 'guest' });
  }
}));