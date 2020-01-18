FROM scratch
ADD server ./
EXPOSE 9000
ENTRYPOINT ["/server"]