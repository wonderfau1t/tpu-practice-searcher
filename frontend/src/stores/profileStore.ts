import { create } from 'zustand';
import axiosInstance, {API_URL} from "../api/instance.ts";
import { ApiResponse } from "../types/apiResponse.ts";

interface CompanyInfo {
  name: string;
  description: string;
  link: string;
}

interface HRsResponse {
  totalCount: number;
  hrs: HR[];
}

interface HR {
  id: number;
  username: string,
  countOfVacancies: number,
}

interface ProfileState {
  companyInfo: CompanyInfo | null;
  hrs: HR[];
  isLoading: boolean;
  error: string | null;
  fetchCompanyInfo: () => Promise<void>;
  fetchCompanyInfoFromVacancy: (id: number) => Promise<void>;
  fetchHRs: () => Promise<void>;
  createHR: (username: string) => Promise<void>;
}

export const useProfileStore = create<ProfileState>((set) => ({
  companyInfo: null,
  hrs: [],
  isLoading: false,
  error: null,

    fetchCompanyInfo: async () => {
      set({ isLoading: true, error: null });
      try {
        const response = await axiosInstance.get<ApiResponse<CompanyInfo>>(`${API_URL}/companies/me`);
        if (response.data.status === 'OK') {
          set({ companyInfo: response.data.result, isLoading: false });
        } else {
          set({ error: 'Server Error: Response status is not OK.', isLoading: false });
        }
      } catch (error) {
        set({ error: 'Server Error: Failed to load company info.', isLoading: false });
      }
    },

    fetchHRs: async () => {
      set({ isLoading: true, error: null});
      try {
        const response = await axiosInstance.get<ApiResponse<HRsResponse>>(`${API_URL}/companies/me/hrs`);
        if (response.data.status === 'OK') {
          set({
            hrs: response.data.result.hrs,
            isLoading: false,
          });
        } else {
          set({ error: 'Server Error: Response status is not OK.', isLoading: false })
        }
      } catch (error) {
        set({ error: 'Server Error: Failed to load hrs.', isLoading: false });
      }
    },

    createHR: async (username: string) => {
      set({ isLoading: true, error: null });
      try {
        const response = await axiosInstance.post<ApiResponse<null>>(`${API_URL}/companies/me/hrs`, { username });
        if (response.data.status === 'OK') {
          set({ isLoading: false });
          await useProfileStore.getState().fetchHRs();
        } else {
          set({ error: 'Server Error: Response status is not OK.', isLoading: false });
        }
      } catch (error) {
        set({ error: 'Server Error: Failed to create HR.', isLoading: false });
      }
    },

    fetchCompanyInfoFromVacancy: async (id: number) => {
      set({ isLoading: true, error: null });
      try {
        const response = await axiosInstance.get<ApiResponse<CompanyInfo>>(`/companies/${id}`);
        if (response.data.status === 'OK') {
          set({ companyInfo: response.data.result, isLoading: false });
        } else {
          set({ error: 'Server Error: Response status is not OK.', isLoading: false });
        }
      } catch (error) {
        set({ error: 'Server Error: Failed to fetch company.', isLoading: false });
      }
    }
}));