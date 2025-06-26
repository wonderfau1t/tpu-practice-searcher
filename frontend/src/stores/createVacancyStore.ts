import { create } from 'zustand';
import axiosInstance from '../api/instance.ts';
import {toast} from "react-toastify";
import { ApiResponse } from "../types/apiResponse.ts";

interface Format {
  id: number;
  name: string;
}

interface PaymentMethod {
  id: number;
  name: string;
}

interface VacancyForm {
  companyName: string;
  vacancyName: string;
  formatId: number | null;
  courseIds: number[];
  keywords: string[];
  deadlineAt: string;
  paymentForAccommodationId: number | null;
  paymentForAccommodationDetails: string;
  farePaymentId: number | null;
  farePaymentDetails: string;
  description: {
    workplace: string;
    position: string;
    salary: number | string;
    requirements: string;
    food: string;
    conditions: string;
    additionalInfo: string;
  };
}

interface VacancyState {
  form: VacancyForm;
  formats: Format[];
  courses: { id: number; name: string }[];
  farePaymentMethods: PaymentMethod[];
  accommodationPaymentMethods: PaymentMethod[];
  isLoading: boolean;
  hasFetchedOptions: boolean;
  hasFetchedFilters: boolean;
  setFormField: <K extends keyof VacancyForm>(field: K, value: VacancyForm[K]) => void;
  setDescriptionField: <K extends keyof VacancyForm['description']>(
    field: K,
    value: VacancyForm['description'][K]
  ) => void;
  resetForm: () => void;
  fetchFormOptions: () => Promise<void>;
  submitVacancy: () => Promise<void>;
  submitVacancyWithoutCompany: () => Promise<void>;
  fetchFilters: () => Promise<void>;
}

export const useVacancyStore = create<VacancyState>((set, get) => ({
  form: {
    companyName: '',
    vacancyName: '',
    formatId: null,
    courseIds: [],
    keywords: [],
    deadlineAt: '',
    paymentForAccommodationId: null,
    paymentForAccommodationDetails: '',
    farePaymentId: null,
    farePaymentDetails: '',
    description: {
      workplace: '',
      position: '',
      salary: '',
      requirements: '',
      food: '',
      conditions: '',
      additionalInfo: '',
    },
  },
  formats: [],
  courses: [],
  farePaymentMethods: [],
  accommodationPaymentMethods: [],
  isLoading: false,
  hasFetchedOptions: false,
  hasFetchedFilters: false,

  setFormField: (field, value) => {
    set((state) => ({
      form: { ...state.form, [field]: value },
    }));
  },

  setDescriptionField: (field, value) => {
    set((state) => ({
      form: {
        ...state.form,
        description: {
          ...state.form.description,
          [field]: value,
        },
      },
    }));
  },

  resetForm: () => {
    set({
      form: {
        companyName: '',
        vacancyName: '',
        formatId: null,
        courseIds: [],
        keywords: [],
        deadlineAt: '',
        paymentForAccommodationId: null,
        paymentForAccommodationDetails: '',
        farePaymentId: null,
        farePaymentDetails: '',
        description: {
          workplace: '',
          position: '',
          salary: '',
          requirements: '',
          food: '',
          conditions: '',
          additionalInfo: '',
        },
      },
    });
  },

  fetchFormOptions: async () => {
    // Проверяем, были ли данные уже загружены
    if (get().hasFetchedOptions) {
      return;
    }

    set({ isLoading: true });

    try {
      const [formatsRes, coursesRes, fareRes, accommodationRes] = await Promise.all([
        axiosInstance.get<ApiResponse<{ formats: Format[] }>>('/references/formats'),
        axiosInstance.get<ApiResponse<{ courses: { id: number; name: string }[] }>>('/references/courses'),
        axiosInstance.get<ApiResponse<{ paymentMethods: PaymentMethod[] }>>('/references/farePaymentMethods'),
        axiosInstance.get<ApiResponse<{ paymentMethods: PaymentMethod[] }>>('/references/accommodationPaymentMethods'),
      ]);

      set({
        formats: formatsRes.data.result.formats,
        courses: coursesRes.data.result.courses,
        farePaymentMethods: fareRes.data.result.paymentMethods,
        accommodationPaymentMethods: accommodationRes.data.result.paymentMethods,
        isLoading: false,
        hasFetchedOptions: true, // Устанавливаем флаг после успешной загрузки
      });
    } catch (error) {
      console.error('Ошибка загрузки данных:', error);
      set({ isLoading: false, hasFetchedOptions: false });
    }
  },

  submitVacancy: async () => {
    set({ isLoading: true });
    try {
      const { form } = get();

      await axiosInstance.post('/vacancies', {
        vacancyName: form.vacancyName,
        formatID: form.formatId,
        courses: form.courseIds,
        keywords: form.keywords,
        deadlineAt: form.deadlineAt,
        paymentForAccommodationID: form.paymentForAccommodationId,
        paymentForAccommodationDetails: form.paymentForAccommodationDetails,
        farePaymentID: form.farePaymentId,
        farePaymentDetails: form.farePaymentDetails,
        description: {
          workplace: form.description.workplace,
          position: form.description.position,
          salary: form.description.salary,
          requirements: form.description.requirements,
          food: form.description.food,
          conditions: form.description.conditions,
          additionalInfo: form.description.additionalInfo,
        },
      });

      try {
        get().resetForm();
      } catch (resetError) {
        console.error('Ошибка в resetForm', resetError);
      }

      set({ isLoading: false });
    } catch (error) {
      set({ isLoading: false });
      throw error;
    }
  },
  submitVacancyWithoutCompany: async () => {
    set({ isLoading: true });
    try {
      const { form } = get();

      await axiosInstance.post(
        "/vacancies",
        {
          companyName: form.companyName,
          vacancyName: form.vacancyName,
          formatID: form.formatId,
          courses: form.courseIds,
          keywords: form.keywords,
          deadlineAt: form.deadlineAt,
          paymentForAccommodationID: form.paymentForAccommodationId,
          paymentForAccommodationDetails: form.paymentForAccommodationDetails,
          farePaymentID: form.farePaymentId,
          farePaymentDetails: form.farePaymentDetails,
          description: {
            workplace: form.description.workplace,
            position: form.description.position,
            salary: form.description.salary,
            requirements: form.description.requirements,
            food: form.description.food,
            conditions: form.description.conditions,
            additionalInfo: form.description.additionalInfo,
          },
        },
      );

      try {
        get().resetForm();
      } catch (resetError) {
        console.error('Ошибка в resetForm', resetError);
      }

      set({ isLoading: false });
    } catch (error) {
      set({ isLoading: false });
      throw error;
    }
  },
  fetchFilters: async () => {
    if (get().hasFetchedFilters) {
      return;
    }

    set({ isLoading: true });

    try {
      const [coursesRes] = await Promise.all([
        axiosInstance.get<ApiResponse<{ courses: { id: number; name: string }[] }>>('/references/courses'),
      ]);

      set({
        courses: coursesRes.data.result.courses,
        isLoading: false,
        hasFetchedFilters: true,
      });
    } catch (error) {
      console.error('Ошибка загрузки категорий и направлений:', error);
      set({ isLoading: false, hasFetchedFilters: false });
      toast.error('Не удалось загрузить фильтры');
    }
  },
}));