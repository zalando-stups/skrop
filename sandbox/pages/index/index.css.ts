import { css } from "emotion";

export const logoClass = css({
  "@media only screen and (max-width: 768px)": {
    marginBottom: "2rem",
    marginTop: "2rem"
  },

  marginBottom: "2rem",
  marginTop: "12rem"
});

export const communityClass = css({
  "@media only screen and (max-width: 768px)": {
    marginTop: "0"
  },

  marginTop: "16rem"
});

export const indexTopClass = css({
  "-webkit-clip-path":
    "polygon(0% 0%,100% 0%,100% 60%,100% 100%,20% 100%,20% 56%,0 56%)",
  "clip-path":
    "polygon(0% 0%,100% 0%,100% 60%,100% 100%,20% 100%,20% 56%,0 56%)"
});
