import {Link} from "react-router-dom";

export default function NotFoundPage() {
  return (
    <>
      <div>404 Not Found</div>
      <Link to="/">Go to Home Page</Link>
    </>
  );
}