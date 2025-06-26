import { create } from 'zustand';
import axiosInstance, {API_URL} from "../api/instance.ts";
import { ApiResponse } from "../types/apiResponse.ts";

// interface VacancyCardResponse {
//   totalCount: number;
//   vacancies: VacancyCard[];
// }

// interface StudentVacancyCardResponse {
//   totalCount: number;
//   vacancies: StudentVacancyCard[];
// }

// interface ModeratorVacancyCardResponse {
//   totalCount: number;
//   vacancies: ModeratorVacancyCard[];
// }

interface VacancyListResponse<T> {
  totalCount: number;
  vacancies: T[];
}

interface StudentVacancyCard {
  id: number;
  name: string;
  companyName: string;
  countOfReplies: number;
  isCreatedByUser: boolean;
}

interface ModeratorVacancyCard {
  id: number;
  name: string;
  companyName: string;
  countOfReplies: number;
  isCreatedByUser: boolean;
}

interface VacancyCard {
  id: number;
  name: string;
  courses: string[];
  countOfReplies: number;
}

interface ReplyResponse {
  status: string;
  result: string;
}

interface TgUserDataInVacancy {
  username: string;
  phoneNumber: string;
}

export interface Course {
  courseId: number;
  name: string;
}

export interface VacancyInfo {
  id: number;
  companyID?: number;
  hrInfo?: TgUserDataInVacancy;
  repliedStudents?: TgUserDataInVacancy[];
  vacancyName: string;
  companyName: string;
  format: string;
  formatID: number | null;
  courses: Course[];
  keywords: string[];
  deadlineAt: string;
  paymentForAccommodation: string;
  paymentForAccommodationID: number | null;
  paymentForAccommodationDetails: string;
  farePayment: string;
  farePaymentID: number | null;
  farePaymentDetails: string;
  description: {
    workplace: string;
    position: string;
    salary: string;
    requirements: string;
    food: string;
    conditions: string;
    additionalInfo: string;
  };
  isReplied?: boolean;
  isCreatedByUser?: boolean;
  hasCompanyProfile?: boolean;
}

interface VacancyById {
  vacancyInfo: VacancyInfo;
}

interface VacancyCardState {
  vacancies: VacancyCard[];
  moderatorVacancies: ModeratorVacancyCard[];
  allVacancies: StudentVacancyCard[];
  repliedVacancies: StudentVacancyCard[];
  totalCount: number;
  currentVacancy: VacancyInfo | null;
  isLoading: boolean;
  error: string | null;
  fetchCompanyVacancies: () => Promise<void>;
  fetchModeratorVacancies: () => Promise<void>;
  fetchAllVacancies: () => Promise<void>;
  fetchRepliedVacancies: () => Promise<void>;
  fetchVacancyById: (id: number) => Promise<void>;
  fetchVacancyByModerator: (id: number) => Promise<void>;
  replyToVacancy: (vacancyID: number) => Promise<void>;
  declineVacancy: (vacancyID: number) => Promise<void>;
  filterVacancies: (params: { courseIDs: number[] }) => Promise<void>;
  searchVacancies: (query: string) => Promise<void>;
  deleteVacancy: (vacancyID: number) => Promise<void>;
  editVacancy: (vacancyID: number, updatedData: Partial<VacancyInfo>) => Promise<void>;
}

export const useVacancyCardStore = create<VacancyCardState>((set) => ({
  vacancies: [],
  moderatorVacancies: [],
  allVacancies: [],
  repliedVacancies: [],
  totalCount: 0,
  currentVacancy: null,
  isLoading: false,
  error: null,
  fetchCompanyVacancies: async () => {
    set ({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.get<ApiResponse<VacancyListResponse<VacancyCard>>>(`${API_URL}/vacancies`);
      if (response.data.status === 'OK') {
        set({
          vacancies: response.data.result.vacancies,
          totalCount: response.data.result.totalCount,
          isLoading: false,
        });
      } else {
        set({
          error: 'Server Error: Response status is not OK.',
          isLoading: false
        });
      }
    } catch (error) {
      set({
        error: 'Server Error: Failed to load vacancies.',
        isLoading: false
      });
    }
  },

  fetchModeratorVacancies: async () => {
    set ({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.get<ApiResponse<VacancyListResponse<ModeratorVacancyCard>>>(
        "/vacancies",
        );
      if (response.data.status === 'OK') {
        set({
          moderatorVacancies: response.data.result.vacancies,
          totalCount: response.data.result.totalCount,
          isLoading: false,
        });
      } else {
        set({
          error: 'Server Error: Response status is not OK.',
          isLoading: false
        });
      }
    } catch (error) {
      set({
        error: 'Server Error: Failed to load vacancies.',
        isLoading: false
      });
    }
  },

  fetchAllVacancies: async () => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.get<ApiResponse<VacancyListResponse<StudentVacancyCard>>>(
        `/vacancies`,
      );
      if (response.data.status === 'OK') {
        set({
          allVacancies: response.data.result.vacancies,
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

  fetchRepliedVacancies: async () => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.get<ApiResponse<VacancyListResponse<StudentVacancyCard>>>(`/vacancies/replies`);
      if (response.data.status === 'OK') {
        set({
          repliedVacancies: response.data.result.vacancies,
          isLoading: false,
        });
      } else {
        set({
          error: 'Server Error: Response status is not OK.',
          isLoading: false
        })
      }
    } catch {
      set({
        error: 'Server Error: Failed to load replied vacancies.',
      })
    }
  },

  fetchVacancyById: async (id: number) => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.get<ApiResponse<VacancyById>>(`/vacancies/${id}`);
      if (response.data.status === 'OK') {
        set({
          currentVacancy: response.data.result.vacancyInfo,
          isLoading: false,
        });
      } else {
        set({
          error: 'Server Error: Response status is not OK.',
          isLoading: false,
        });
      }
    } catch (error) {
      set({
        error: 'Server Error: Failed to load vacancy.',
        isLoading: false,
      });
    }
  },

  fetchVacancyByModerator: async (id: number) => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.get<ApiResponse<VacancyById>>(
        `/vacancies/${id}`,
      );
      if (response.data.status === 'OK') {
        set({
          currentVacancy: response.data.result.vacancyInfo,
          isLoading: false,
        });
      } else {
        set({
          error: 'Server Error: Response status is not OK.',
          isLoading: false,
        });
      }
    } catch (error) {
      set({
        error: 'Server Error: Failed to load vacancy.',
        isLoading: false,
      });
    }
  },

  replyToVacancy: async (vacancyID: number) => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.post<ApiResponse<ReplyResponse>>(
        `/vacancies/${vacancyID}/replies`,
      );

      if (response.data.status === 'OK') {
        set((state) => ({
          isLoading: false,
          error: null,
          currentVacancy: state.currentVacancy
            ? { ...state.currentVacancy, isReplied: true }
            : null,
        }));
      }
    } catch (error) {
      set({
        error: null,
        isLoading: false,
      });
      throw error;
    }
  },

  declineVacancy: async (vacancyID: number) => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.delete<ApiResponse<ReplyResponse>>(
        `/vacancies/${vacancyID}/replies`,
      );

      if (response.data.status === 'OK') {
        set((state) => ({
          isLoading: false,
          error: null,
          currentVacancy: state.currentVacancy
            ? { ...state.currentVacancy, isReplied: false }
            : null,
        }));
      } else {
        throw new Error('Вы уже откликнулись или сервер вернул ошибку');
      }
    } catch (error) {
      set({
        error: null,
        isLoading: false,
      });
      throw error;
    }
  },

  filterVacancies: async (params: { courseIDs?: number[] }) => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.get<ApiResponse<VacancyListResponse<StudentVacancyCard>>>(
        `/vacancies/filter`,
        {
          params: {
            course_ids: params.courseIDs?.join(','),
          }
        }
      );
      if (response.data.status === 'OK') {
        set({
          allVacancies: response.data.result.vacancies,
          totalCount: response.data.result.totalCount,
          isLoading: false
        });
      } else {
        set({
          error: 'Server Error: Response status is not OK.',
          isLoading: false
        });
      }
    } catch {
      set({
        error: 'Server Error: Failed to load filtered vacancies.',
        isLoading: false
      });
    }
  },

  searchVacancies: async (query: string) => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.get<ApiResponse<VacancyListResponse<StudentVacancyCard>>>(`/vacancies/search`,
        {
          params: { query }
        }
      );
      if (response.data.status === 'OK') {
        set({
          allVacancies: response.data.result.vacancies,
          totalCount: response.data.result.totalCount,
          isLoading: false,
        })
      } else {
        set({
          error: 'Server Error: Response status is not OK.',
          isLoading: false,
        })
      }
    } catch (error) {
      set({ error: 'Server Error: Failed to search vacancies.', isLoading: false });
    }
  },

  deleteVacancy: async (vacancyID: number) => {
    set({ isLoading: true, error: null });
    try {
      const response = await axiosInstance.delete<ApiResponse<VacancyById>>(
        `/vacancies/${vacancyID}`,
      );

      if (response.data.status === 'OK') {
        set((state) => ({
          isLoading: false,
          error: null,
          currentVacancy: null,
          vacancies: state.vacancies.filter(v => v.id !== vacancyID),
          moderatorVacancies: state.moderatorVacancies.filter(v => v.id !== vacancyID),
        }));
      } else {
        throw new Error('Вы уже откликнулись или сервер вернул ошибку');
      }
    } catch (error) {
      set({
        error: null,
        isLoading: false,
      });
      throw error;
    }
  },
  editVacancy: async (vacancyId: number, updatedData: Partial<VacancyInfo>) => {
    set({ isLoading: true, error: null });

    try {
      const response = await axiosInstance.put(`/vacancies/${vacancyId}`,
        updatedData,
      );
      if (response.data.status === 'OK') {
        set({
          isLoading: false,
          error: null });
      } else {
        set({ error: 'Server Error: Response status is not OK.', isLoading: false });
      }
    } catch (error) {
      set({ error: 'Server Error: Failed to reg company.', isLoading: false });
    }
  },
}))

