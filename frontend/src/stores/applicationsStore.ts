import { create } from 'zustand';
import {ApiResponse} from "../types/apiResponse.ts";
import axiosInstance from "../api/instance.ts";

interface CompanyApplication {
  companyID: number,
  name: string,
  description: string,
  link: string,
  registeredAt: string,
  HRUsername: string,
}

interface CompanyApplicationsState {
  companies: CompanyApplication[];
  selectedApplication: CompanyApplication | null;
  totalCount: number;
  isLoading: boolean;
  error: string | null;
  fetchCompaniesApplications: () => Promise<void>;
  fetchCompanyApplicationByID: (companyID: number) => Promise<void>;
  approveCompanyApplication: (companyID: number) => Promise<void>;
  rejectCompanyApplication: (companyID: number, message: string) => Promise<void>;
  clear: () => void;
}

export const useApplicationsStore = create<CompanyApplicationsState>((set) => ({
  companies: [],
  selectedApplication: null,
  totalCount: 0,
  isLoading: false,
  error: null,

  fetchCompaniesApplications: async () => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.get<ApiResponse<CompanyApplicationsState>>('/companies/requests');
      if (response.data.status === 'OK') {
        set({
          companies: response.data.result.companies,
          totalCount: response.data.result.totalCount,
          isLoading: false
        });
      } else {
        set({
          error: 'Server Error: Response status is not OK.',
          isLoading: false
        })
      }
    } catch {
      set({
        error: 'Server Error: Failed to load vacancies.',
      })
    }
  },

  fetchCompanyApplicationByID: async (companyID: number) => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.get<ApiResponse<CompanyApplication>>(`/companies/requests/${companyID}`,
        {
          headers: {
            "Authorization": `Bearer ${localStorage.getItem("accessToken")}`,
          }
        }
      );
      if (response.data.status === 'OK') {
        set({
          selectedApplication: response.data.result,
          isLoading: false
        });
      } else {
        set({
          error: 'Server Error: Response status is not OK.',
          isLoading: false
        })
      }
    } catch {
      set({
        error: 'Server Error: Failed to load vacancies.',
      })
    }
  },

  clear: () => {
    set({
      companies: [],
      totalCount: 0,
      error: null,
    });
  },
  approveCompanyApplication: async (companyID: number) => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.patch(
        `/companies/requests/${companyID}/accept`,
      );

      if (response.data.status === 'OK') {
        set({
          isLoading: false
        });
      } else {
        set({
          error: 'Ошибка сервера при подтверждении заявки',
          isLoading: false
        });
      }
    } catch {
      set({
        error: 'Не удалось подтвердить заявку',
        isLoading: false
      });
    }
  },

  rejectCompanyApplication: async (companyID: number, message: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.patch(
        `/companies/requests/${companyID}/reject`,
        {
          message
        },
      );

      if (response.data.status === 'OK') {
        set({
          isLoading: false
        });
      } else {
        set({
          error: 'Ошибка сервера при отклонении заявки',
          isLoading: false
        });
      }
    } catch {
      set({
        error: 'Не удалось отклонить заявку',
        isLoading: false
      });
    }
  },

}));