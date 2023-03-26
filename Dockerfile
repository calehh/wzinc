FROM public.ecr.aws/zinclabs/zincobserve:latest
EXPOSE 4080
ENTRYPOINT ["/go/bin/zincsearch"]