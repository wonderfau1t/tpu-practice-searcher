import React, { useState, useEffect } from "react";
import { GridLoader } from "react-spinners";
import "./Loader.scss";

interface LoaderProps {
  isLoading: boolean;
  delay?: number; // Задержка в миллисекундах перед показом лоадера
}

const Loader: React.FC<LoaderProps> = ({ isLoading, delay = 300 }) => {
  const [showLoader, setShowLoader] = useState(false);

  useEffect(() => {
    let timer: NodeJS.Timeout | null = null;

    if (isLoading) {
      // Показываем лоадер только после задержки
      timer = setTimeout(() => {
        setShowLoader(true);
      }, delay);
    } else {
      // Скрываем лоадер сразу, если загрузка завершена
      setShowLoader(false);
    }

    return () => {
      if (timer) clearTimeout(timer);
    };
  }, [isLoading, delay]);

  if (!showLoader) return null;

  return (
    <div className="loader-container">
      <GridLoader
        color="#28be46" // Зеленый цвет, как ты хотел ранее
        size={20} // Размер лоадера (можно настроить)
        margin={4} // Отступ между точками
        loading={showLoader}
        aria-label="Loading Spinner"
      />
    </div>
  );
};

export default Loader;