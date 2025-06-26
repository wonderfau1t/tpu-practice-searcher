import {tg} from "../lib/telegram.ts";
import {toast} from "react-toastify";

const handleRequestContact = () => {
  tg.requestContact((wasShared: boolean)=> {
    if (wasShared) {
      toast.success('Вы успешно поделились номером!');
    }
  });
}

export default handleRequestContact;