export function getMsg(error) {
  if (error == null) return "";
  if (typeof error === "string") return error;
  if (error.message != undefined) return error.message;
  console.warn("can't get message from error:", error);
  return error;
}
