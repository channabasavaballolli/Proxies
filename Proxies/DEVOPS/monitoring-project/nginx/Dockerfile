#uses official ngix image
FROM nginx:latest  
#replace default web page
COPY index.html /usr/share/nginx/html/index.html 
#//replace conffig
COPY nginx.conf /etc/nginx/nginx.conf 
#//contains listens on 80 port
EXPOSE 80 
# Start NGINX
CMD ["nginx","-g","daemon off;"]