import { create } from "zustand";
import {tg} from "../lib/telegram.ts";
import axios from "axios";
import axiosInstance from "../api/instance.ts";
import {ApiResponse} from "../types/apiResponse.ts";

interface Company {
  name: string;
  description: string;
  link: string;
}

interface CompanyState {
  company: Company | null;
  isLoading: boolean;
  error: string | null;
  registerCompany: (company: Company) => Promise<void>;
  editCompany: (company: Company) => Promise<void>;
}

export const useCompanyStore = create<CompanyState>((set) => ({
  company: null,
  isLoading: false,
  error: null,

  registerCompany: async (company: Company) => {
    set({ isLoading: true, error: null });

    try {
      const response = await axios.post(
        "https://api.tpupractice.ru/companies",
        {
          name: company.name,
          description: company.description,
          link: company.link,
        },
        {
          headers: {
            "Authorization": `tma ${tg.initData}`
          }
        }
      );
      if (response.data.status === 'OK') {
        set({ isLoading: false, error: null });
      } else {
        set({ error: 'Server Error: Response status is not OK.', isLoading: false });
      }
    } catch (error) {
      set({ error: 'Server Error: Failed to reg company.', isLoading: false });
    }
  },

  editCompany: async (company: { name: string, description: string, link: string }) => {
    set({ isLoading: true, error: null });

    try {
      const response = await axiosInstance.put<ApiResponse<Company>>(
        `/companies/me`,
        {
          name: company.name,
          description: company.description,
          link: company.link,
        },
      );
      if (response.data.status === 'OK') {
        set({ isLoading: false, error: null });
      } else {
        set({ error: 'Server Error: Response status is not OK.', isLoading: false });
      }
    } catch (error) {
      set({ error: 'Server Error: Failed to reg company.', isLoading: false });
    }
  },
}));